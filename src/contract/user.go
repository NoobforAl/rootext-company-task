package contract

import (
	"context"
	"ratblog/internal/entity"
)

type User interface {
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	GetAllUsersWithPagination(ctx context.Context, page, limit int) (*entity.SliceOfUser, error)
	GetUserByID(ctx context.Context, id int) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id int) error
}
