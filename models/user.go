package models

// JSON representation of user
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// JSON representation of user without password, for API response
type UserResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// JSON representation for created JWT tokens
type JWTTokenResponse struct {
	Token string `json:"token"`
}
