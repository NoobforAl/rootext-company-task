package repository

import (
	"context"
	"ratblog/internal/database"
	"ratblog/internal/entity"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

func (r *repository) GetAllPosts(ctx context.Context) ([]*entity.Post, error) {
	posts, err := r.db.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}

	var result []*entity.Post
	for _, post := range posts {
		result = append(result, &entity.Post{
			ID:        post.ID,
			UserID:    post.UserID,
			Title:     post.Title,
			Content:   post.Content,
			RateScore: post.TotalRating,
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}
	return result, nil
}

func (r *repository) GetPostByID(ctx context.Context, id int) (*entity.Post, error) {
	post, err := r.db.GetPostByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &entity.Post{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		RateScore: post.TotalRating,
		CreatedAt: post.CreatedAt.Time,
		UpdatedAt: post.UpdatedAt.Time,
	}, nil
}

func (r *repository) GetAllPostsWithPagination(ctx context.Context, page, limit int) (*entity.SliceOfPost, error) {
	offset := (page - 1) * limit
	params := database.GetAllPostsWithPaginationParams{
		Column1: float64(limit),
		Offset:  int64(offset),
		Limit:   int64(limit),
	}

	posts, err := r.db.GetAllPostsWithPagination(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*entity.Post
	for _, post := range posts {
		result = append(result, &entity.Post{
			ID:        post.ID,
			UserID:    post.UserID,
			Title:     post.Title,
			Content:   post.Content,
			RateScore: post.TotalRating,
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}

	total := 0
	if len(result) > 0 {
		total = int(posts[0].TotalCount)
	}

	return &entity.SliceOfPost{
		Posts:    result,
		Page:     page,
		Size:     limit,
		MaxPosts: total,
	}, nil
}

func (r *repository) GetTopPostsInPeriodWithPagination(ctx context.Context, period string, page, limit int) (*entity.SliceOfPost, error) {
	offset := (page - 1) * limit

	// pars period string 1h or 1m or 1s to microseconds
	timePeriod, _ := time.ParseDuration(period)
	interval := pgtype.Interval{
		Microseconds: int64(timePeriod / time.Microsecond),
		Valid:        true,
	}

	params := database.GetTopPostsInPeriodWithPaginationParams{
		Offset:  int64(offset),
		Limit:   float64(limit),
		Column1: interval,
	}

	posts, err := r.db.GetTopPostsInPeriodWithPagination(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*entity.Post
	for _, post := range posts {
		result = append(result, &entity.Post{
			ID:        post.ID,
			UserID:    post.UserID,
			Title:     post.Title,
			Content:   post.Content,
			RateScore: post.TotalRating.Int32,
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}

	total := 0
	if len(result) > 0 {
		total = int(posts[0].TotalCount)
	}

	return &entity.SliceOfPost{
		Posts:    result,
		Page:     page,
		Size:     limit,
		MaxPosts: total,
	}, nil
}

func (r *repository) CreatePost(ctx context.Context, post *entity.Post) error {
	_, err := r.db.CreatePost(ctx, database.CreatePostParams{
		UserID:  int32(post.UserID),
		Title:   post.Title,
		Content: post.Content,
	})
	return err
}

func (r *repository) UpdatePost(ctx context.Context, post *entity.Post) error {
	err := r.db.UpdatePost(ctx, database.UpdatePostParams{
		ID:      int32(post.ID),
		UserID:  int32(post.UserID),
		Title:   post.Title,
		Content: post.Content,
	})
	return err
}

func (r *repository) DeletePost(ctx context.Context, id, userId int) error {
	err := r.db.DeletePost(ctx, database.DeletePostParams{
		ID:     int32(id),
		UserID: int32(userId),
	})
	return err
}

func (r *repository) GetAllUserRatedPostsWithPagination(ctx context.Context, userId, page, limit int) (*entity.SliceOfPost, error) {
	offset := (page - 1) * limit
	params := database.GetAllUserRatedPostsWithPaginationParams{
		UserID: int32(userId),
		Offset: int64(offset),
		Limit:  float64(limit),
	}

	posts, err := r.db.GetAllUserRatedPostsWithPagination(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*entity.Post
	for _, post := range posts {
		result = append(result, &entity.Post{
			ID:        post.ID,
			UserID:    post.UserID,
			Title:     post.Title,
			Content:   post.Content,
			RateScore: post.RatingValue,
			CreatedAt: post.CreatedAt.Time,
			UpdatedAt: post.UpdatedAt.Time,
		})
	}

	total := 0
	if len(result) > 0 {
		total = int(posts[0].TotalCount)
	}

	return &entity.SliceOfPost{
		Posts:    result,
		Page:     page,
		Size:     limit,
		MaxPosts: total,
	}, nil
}

func (r *repository) RateUpPost(ctx context.Context, id, userId int) error {
	return r.db.UpsertRatePost(ctx, database.UpsertRatePostParams{
		PostID:      int32(id),
		UserID:      int32(userId),
		RatingValue: 1,
	})
}

func (r *repository) RateDownPost(ctx context.Context, id, userId int) error {
	return r.db.UpsertRatePost(ctx, database.UpsertRatePostParams{
		PostID:      int32(id),
		UserID:      int32(userId),
		RatingValue: -1,
	})
}

func (r *repository) DeleteRatePost(ctx context.Context, id, userId int) error {
	return r.db.DeleteRatePost(ctx, database.DeleteRatePostParams{
		PostID: int32(id),
		UserID: int32(userId),
	})
}
