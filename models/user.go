package models

// User is the JSON representation for users over the REST API
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse is the JSON representation of user without password, for API response
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// JWTTokenResponse is the JSON representation for created JWT tokens
type JWTTokenResponse struct {
	Token string `json:"token"`
}
