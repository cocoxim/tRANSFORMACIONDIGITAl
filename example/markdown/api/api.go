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
//	@Produce	json
//	@Success	200	{array}	api.UsersCollection	"ok"
//	@Router		/admin/user/ [get]
func ListUsers(w http.ResponseWriter, r *http.Request) {
	// write your code
}

// GetUser example
//
//	@Summary	Read user from the store
//	@Tags		admin
//	@Accept		json
//	@Produce	json
//	@Param		id	path		int	true	"User Id"
//	@Success	200	{object}	api.User
//	@Failure	400	{object}	api.APIError	"We need ID!!"
//	@Failure	404	{object}	api.APIError	"Can not find ID"
//	@Router		/admin/user/{id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	// write your code
}

// AddUser example
//
//	@Summary	Add a new user to the