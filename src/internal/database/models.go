// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Post struct {
	ID        int32
	UserID    int32
	Title     string
	Content   string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type PostRatingSummary struct {
	PostID      int32
	TotalRating pgtype.Int4
}

type PostTag struct {
	PostID int32
	TagID  int32
}

type Rating struct {
	ID          int32
	PostID      int32
	UserID      int32
	RatingValue int32
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
}

type Tag struct {
	ID        int32
	Name      string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}

type UserInfo struct {
	ID        int32
	Email     string
	Username  string
	FirstName string
	LastName  string
	Role      string
	Password  string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}
