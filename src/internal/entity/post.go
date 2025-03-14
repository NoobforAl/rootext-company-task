package entity

import (
	"time"
)

type Post struct {
	ID        int32
	UserID    int32
	Title     string
	Content   string
	RateScore int32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SliceOfPost struct {
	Posts    []*Post
	Page     int
	Size     int
	MaxPosts int
}
