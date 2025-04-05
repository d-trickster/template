package db

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
	ErrInternal = errors.New("db error")
)

type DB interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
}

type User struct {
	ID string
}
