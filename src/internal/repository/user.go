package repository

import (
	"context"
	"ratblog/database"
	"ratblog/internal/entity"
)

func (r *repository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	users, err := r.db.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var result []*entity.User
	for _, user := range users {
		result = append(result, &entity.User{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
		})
	}
	return result, nil
}

func (r *repository) GetUserByID(ctx context.Context, id int) (*entity.User, error) {
	user, err := r.db.GetUserByID(ctx, int32(id))
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
	}, nil
}

func (r *repository) GetAllUsersWithPagination(ctx context.Context, page, limit int) (*entity.SliceOfUser, error) {
	offset := (page - 1) * limit
	params := database.GetAllUsersWithPaginationParams{
		Offset: int64(offset),
		Limit:  float64(limit),
	}

	rows, err := r.db.GetAllUsersWithPagination(ctx, params)
	if err != nil {
		return nil, err
	}

	var users []*entity.User
	for _, row := range rows {
		users = append(users, &entity.User{
			ID:        row.ID,
			Email:     row.Email,
			Username:  row.Username,
			FirstName: row.FirstName,
			LastName:  row.LastName,
			Role:      row.Role,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}

	maxUsers := 0
	if len(rows) != 0 {
		maxUsers = int(rows[0].TotalCount)
	}

	return &entity.SliceOfUser{
		Users:    users,
		Page:     page,
		Size:     limit,
		MaxUsers: maxUsers,
	}, nil
}

func (r *repository) CreateUser(ctx context.Context, user *entity.User) error {
	params := database.CreateUserParams{
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Role:      user.Role,
		Digest:    user.Password,
	}
	_, err := r.db.CreateUser(ctx, params)
	return err
}

func (r *repository) UpdateUser(ctx context.Context, user *entity.User) error {
	params := database.UpdateUserParams{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	return r.db.UpdateUser(ctx, params)
}

func (r *repository) DeleteUser(ctx context.Context, id int) error {
	return r.db.DeleteUser(ctx, int32(id))
}
