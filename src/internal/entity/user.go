package entity

import "time"

type User struct {
	ID        int32
	Email     string
	Username  string
	FirstName string
	LastName  string
	Role      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SliceOfUser struct {
	Users    []*User
	Page     int
	Size     int
	MaxUsers int
}
