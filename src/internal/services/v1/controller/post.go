package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"ratblog/internal/entity"
	"ratblog/internal/services/v1/schema"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func (con *Controller) GetPosts() echo.HandlerFunc {
	return func(c echo.Context) error {
		pageStr := c.QueryParam("page")
		sizeStr := c.QueryParam("size")

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			size = 10
		}

		ctx := c.Request().Context()
		slicPosts, err := con.repo.GetAllPostsWithPagination(ctx, page, size)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		schemaPosts := make([]schema.PostInfo, 0, 16)
		for _, post := range slicPosts.Posts {
			schemaPosts = append(schemaPosts, schema.PostInfo{
				ID:        post.ID,
				Title:     post.Title,
				Content:   post.Content,
				CreatedAt: post.CreatedAt,
				UpdatedAt: post.UpdatedAt,
			})
		}

		return c.JSON(http.StatusOK, schema.SliceOfPostInfo{
			Posts: schemaPosts,
			Page:  int32(slicPosts.Page),
			Size:  int32(slicPosts.Size),
			Total: int32(slicPosts.MaxPosts),
		})
	}
}

func (con *Controller) checkCachedData(ctx context.Context, key, countKey string) (needCache, beforeCached bool, err error) {
	res, err := con.redisClient.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, false, err
	}

	if res != "" {
		return false, true, nil
	}

	count, _ := con.redisClient.Incr(ctx, countKey).Result()
	if count == 1 {
		con.redisClient.Expire(ctx, countKey, 5*time.Minute)
	}

	if count > 5 {
		return true, false, nil
	}

	return false, false, nil
}

func (con *Controller) GetPostsByFilter() echo.HandlerFunc {
	return func(c echo.Context) error {
		pageStr := c.QueryParam("page")
		sizeStr := c.QueryParam("size")
		timeInterval := c.QueryParam("timeInterval")

		cacheKey := fmt.Sprintf("posts:cache:%s-%s-%s", timeInterval, pageStr, sizeStr)
		countKey := fmt.Sprintf("posts:view:%s-%s-%s", timeInterval, pageStr, sizeStr)

		ctx := c.Request().Context()
		needCache, beforeCached, err := con.checkCachedData(ctx, cacheKey, countKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
		}

		if beforeCached {
			data, err := con.redisClient.Get(c.Request().Context(), cacheKey).Result()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
			}

			var res schema.SliceOfPostInfo
			if err := json.Unmarshal([]byte(data), &res); err != nil {
				return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
			}

			return c.JSON(http.StatusOK, res)
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil {
			page = 1
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			size = 10
		}

		if timeInterval == "" {
			timeInterval = "1h"
		}

		// check if timeInterval is valid
		if _, err := time.ParseDuration(timeInterval); err != nil {
			return c.JSON(http.StatusBadRequest, response{Message: err.Error()})
		}

		slicPosts, err := con.repo.GetTopPostsInPeriodWithPagination(ctx, timeInterval, page, size)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		schemaPosts := make([]schema.PostInfo, 0, 16)
		for _, post := range slicPosts.Posts {
			schemaPosts = append(schemaPosts, schema.PostInfo{
				ID:        post.ID,
				Title:     post.Title,
				Content:   post.Content,
				CreatedAt: post.CreatedAt,
				UpdatedAt: post.UpdatedAt,
			})
		}

		res := schema.SliceOfPostInfo{
			Posts: schemaPosts,
			Page:  int32(slicPosts.Page),
			Size:  int32(slicPosts.Size),
			Total: int32(slicPosts.MaxPosts),
		}

		if needCache {
			resBytes, err := json.Marshal(res)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
			}

			err = con.redisClient.Set(ctx, cacheKey, resBytes, 5*time.Minute).Err()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
			}
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (con *Controller) CreatePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := validateUser(c)
		if err != nil {
			return c.JSON(401, unauthorizedResponse)
		}

		var newPost schema.PostCreate
		if err := c.Bind(&newPost); err != nil {
			return c.JSON(400, badRequestResponse)
		}

		err = newPost.Validate()
		if err != nil {
			return c.JSON(400, response{Message: err.Error()})
		}

		ctx := c.Request().Context()
		postEntity := &entity.Post{
			Title:   newPost.Title,
			Content: newPost.Content,
			UserID:  int32(user.ID),
		}

		err = con.repo.CreatePost(ctx, postEntity)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		return c.JSON(200, response{Message: "Post created successfully"})
	}
}

func (con *Controller) GetPost() echo.HandlerFunc {
	return func(c echo.Context) error {
		postIdStr := c.Param("id")
		postId, err := strconv.Atoi(postIdStr)
		if err != nil {
			return c.JSON(400, badRequestResponse)
		}

		cacheKey := fmt.Sprintf("post:cache:%d", postId)
		countKey := fmt.Sprintf("post:view:%d", postId)

		ctx := c.Request().Context()
		needCache, beforeCached, err := con.checkCachedData(ctx, cacheKey, countKey)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
		}

		if beforeCached {
			data, err := con.redisClient.Get(ctx, cacheKey).Result()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
			}

			var res schema.PostInfo
			if err := json.Unmarshal([]byte(data), &res); err != nil {
				return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
			}

			return c.JSON(http.StatusOK, res)
		}

		post, err := con.repo.GetPostByID(ctx, postId)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		res := schema.PostInfo{
			ID:        post.ID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
		}

		if needCache {
			resBytes, err := json.Marshal(res)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
			}

			err = con.redisClient.Set(ctx, cacheKey, resBytes, 5*time.Minute).Err()
			if err != nil {
				return c.JSON(http.StatusInternalServerError, response{Message: err.Error()})
			}
		}

		return c.JSON(200, res)
	}
}

func (con *Controller) UpdatePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := validateUser(c)
		if err != nil {
			return c.JSON(401, unauthorizedResponse)
		}

		postIdStr := c.Param("id")
		postId, err := strconv.Atoi(postIdStr)
		if err != nil {
			return c.JSON(400, badRequestResponse)
		}

		var updatePost schema.PostInfo
		if err := c.Bind(&updatePost); err != nil {
			return c.JSON(400, badRequestResponse)
		}

		ctx := c.Request().Context()
		postEntity := &entity.Post{
			ID:      int32(postId),
			Title:   updatePost.Title,
			Content: updatePost.Content,
			UserID:  int32(user.ID),
		}

		err = con.repo.UpdatePost(ctx, postEntity)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		return c.JSON(200, response{Message: "Post updated successfully"})
	}
}

func (con *Controller) DeletePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := validateUser(c)
		if err != nil {
			return c.JSON(401, unauthorizedResponse)
		}

		fmt.Println(user)

		postIdStr := c.Param("id")
		postId, err := strconv.Atoi(postIdStr)
		if err != nil {
			return c.JSON(400, badRequestResponse)
		}

		ctx := c.Request().Context()
		err = con.repo.DeletePost(ctx, postId, user.ID)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		return c.JSON(200, response{Message: "Post deleted successfully"})
	}
}

func (con *Controller) UpVotePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := validateUser(c)
		if err != nil {
			return c.JSON(401, unauthorizedResponse)
		}

		postIdStr := c.Param("id")
		postId, err := strconv.Atoi(postIdStr)
		if err != nil {
			return c.JSON(400, badRequestResponse)
		}

		ctx := c.Request().Context()
		err = con.repo.RateUpPost(ctx, postId, user.ID)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		return c.JSON(200, response{Message: "Post upvoted successfully"})
	}
}

func (con *Controller) DownVotePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := validateUser(c)
		if err != nil {
			return c.JSON(401, unauthorizedResponse)
		}

		postIdStr := c.Param("id")
		postId, err := strconv.Atoi(postIdStr)
		if err != nil {
			return c.JSON(400, badRequestResponse)
		}

		ctx := c.Request().Context()
		err = con.repo.RateDownPost(ctx, postId, user.ID)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		return c.JSON(200, response{Message: "Post downvoted successfully"})
	}
}

func (con *Controller) RemoveRatePost() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := validateUser(c)
		if err != nil {
			return c.JSON(401, unauthorizedResponse)
		}

		postIdStr := c.Param("id")
		postId, err := strconv.Atoi(postIdStr)
		if err != nil {
			return c.JSON(400, badRequestResponse)
		}

		ctx := c.Request().Context()
		err = con.repo.DeleteRatePost(ctx, postId, user.ID)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		return c.JSON(200, response{Message: "Post rate removed successfully"})
	}
}
