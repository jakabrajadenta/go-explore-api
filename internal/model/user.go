package model

import "time"

type User struct {
	ID        int64
	Username  string
	Email     string
	FullName  string
	Phone     string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
