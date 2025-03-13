package contract

import (
	"context"
	"ratblog/internal/entity"
)

type Tag interface {
	GetAllTags(ctx context.Context) ([]*entity.Tag, error)
	GetAllTagsWithPagination(ctx context.Context, page, limit int) (*entity.SliceOfTag, error)
	GetTagByID(ctx context.Context, id int) (*entity.Tag, error)
	CreateTag(ctx context.Context, tag *entity.Tag) error
	UpdateTag(ctx context.Context, tag *entity.Tag) error
	DeleteTag(ctx context.Context, id int) error
}
