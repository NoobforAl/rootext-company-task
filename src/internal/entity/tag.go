package entity

import "time"

type Tag struct {
	ID        int32
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SliceOfTag struct {
	Tags    []*Tag
	Page    int
	Size    int
	MaxTags int
}
