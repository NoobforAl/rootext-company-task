package repository

import (
	"context"
	"ratblog/internal/database"
	"ratblog/internal/entity"
)

func (r *repository) GetAllTags(ctx context.Context) ([]*entity.Tag, error) {
	tags, err := r.db.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}

	var result []*entity.Tag
	for _, tag := range tags {
		result = append(result, &entity.Tag{
			ID:        tag.ID,
			Name:      tag.Name,
			CreatedAt: tag.CreatedAt.Time,
			UpdatedAt: tag.UpdatedAt.Time,
		})
	}
	return result, nil
}

func (r *repository) GetAllTagsWithPagination(ctx context.Context, page, limit int) (*entity.SliceOfTag, error) {
	offset := (page - 1) * limit
	params := database.GetAllTagsWithPaginationParams{
		Offset: int64(offset),
		Limit:  float64(limit),
	}

	rows, err := r.db.GetAllTagsWithPagination(ctx, params)
	if err != nil {
		return nil, err
	}

	var tags []*entity.Tag
	for _, row := range rows {
		tags = append(tags, &entity.Tag{
			ID:        row.ID,
			Name:      row.Name,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}

	maxTags := 0
	if len(rows) != 0 {
		maxTags = int(rows[0].TotalCount)
	}

	return &entity.SliceOfTag{
		Tags:    tags,
		Page:    page,
		Size:    limit,
		MaxTags: maxTags,
	}, nil
}

func (r *repository) GetTagByID(ctx context.Context, id int) (*entity.Tag, error) {
	tag, err := r.db.GetTagByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &entity.Tag{
		ID:        tag.ID,
		Name:      tag.Name,
		CreatedAt: tag.CreatedAt.Time,
		UpdatedAt: tag.UpdatedAt.Time,
	}, nil
}

func (r *repository) CreateTag(ctx context.Context, tag *entity.Tag) error {
	_, err := r.db.CreateTag(ctx, tag.Name)
	return err
}

func (r *repository) UpdateTag(ctx context.Context, tag *entity.Tag) error {
	params := database.UpdateTagParams{
		ID:   tag.ID,
		Name: tag.Name,
	}
	return r.db.UpdateTag(ctx, params)
}

func (r *repository) DeleteTag(ctx context.Context, id int) error {
	return r.db.DeleteTag(ctx, int32(id))
}
