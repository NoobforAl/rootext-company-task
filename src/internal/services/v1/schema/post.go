package schema

import "time"

type PostInfo struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	RateScore int32     `json:"rate_score"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SliceOfPostInfo struct {
	Posts []PostInfo `json:"posts"`
	Page  int32      `json:"page"`
	Size  int32      `json:"size"`
	Total int32      `json:"total"`
}

type PostCreate struct {
	Title   string `json:"title" validate:"required,min=5"`
	Content string `json:"content" validate:"required,min=5"`
}

func (p *PostCreate) Validate() error {
	return validate.Struct(p)
}

type PostUpdate struct {
	Title   string `json:"title" validate:"required,min=5"`
	Content string `json:"content" validate:"required,min=5"`
}

func (p *PostUpdate) Validate() error {
	return validate.Struct(p)
}
