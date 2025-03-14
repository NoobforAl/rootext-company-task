package schema

import "time"

type Login struct {
	UserName string `json:"username" validate:"required,min=5"`
	Password string `json:"password" validate:"required,min=8"`
}

func (l *Login) Validate() error {
	return validate.Struct(l)
}

type Register struct {
	UserName  string `json:"username" validate:"required,min=5"`
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email,min=5"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (r *Register) Validate() error {
	return validate.Struct(r)
}

type UserUpdate struct {
	FirstName string `json:"first_name" validate:"required,min=2"`
	LastName  string `json:"last_name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email,min=5"`
	Password  string `json:"password" validate:"required,min=8"`
}

func (u *UserUpdate) Validate() error {
	return validate.Struct(u)
}

type UserInfo struct {
	Email     string    `json:"email"`
	UserName  string    `json:"username"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
