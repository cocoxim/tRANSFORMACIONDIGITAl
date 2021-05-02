package api

import (
	"net/http"
	"time"
)

// User example
type User struct {
	ID       int64
	Email    string
	Password string
}

// UsersCollection example
type UsersCollection []User

// Error example
type APIError struct {
	ErrorCode    int
	ErrorMessage string
	CreatedAt    time.Time
}

// ListUsers example
//
//	@Summary	List users from the store
//	@Tags		admin
//	@Accept		json
//	@P