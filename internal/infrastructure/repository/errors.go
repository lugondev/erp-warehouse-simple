package repository

import "errors"

var (
	ErrRecordNotFound    = errors.New("record not found")
	ErrInsufficientStock = errors.New("insufficient stock quantity")
	ErrDuplicateEntry    = errors.New("duplicate entry")
	ErrInvalidData       = errors.New("invalid data")
	ErrRoleInUse         = errors.New("role is in use by users")
)
