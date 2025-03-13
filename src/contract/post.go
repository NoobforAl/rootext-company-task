package contract

import (
	"context"
	"ratblog/internal/entity"
)

type Post interface {
	GetAllPosts(ctx context.Context) ([]*entity.Post, error)
	GetPostByID(ctx context.Context, id int) (*entity.Post, error)
	GetAllPostsWithPagination(ctx context.Context, page, limit int) (*entity.SliceOfPost, error)
	GetTopPostsInPeriodWithPagination(ctx context.Context, period string, page, limit int) (*entity.SliceOfPost, error)
	CreatePost(ctx context.Context, post *entity.Post) error
	UpdatePost(ctx context.Context, post *entity.Post) error
	DeletePost(ctx context.Context, id, userId int) error

	GetAllUserRatedPostsWithPagination(ctx context.Context, userId, page, limit int) (*entity.SliceOfPost, error)
	RateUpPost(ctx context.Context, id, userId int) error
	RateDownPost(ctx context.Context, id, userId int) error
	DeleteRatePost(ctx context.Context, id, userId int) error
}
