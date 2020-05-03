package models

// Game is the representation of a game
// in the database.
type Game struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// User is the representation of a user
// in the database.
type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	HashedPassword []byte
}

// UserResponse is the response that is returned
// when a user is logged in.
type UserResponse struct {
	Token string `json:"token"`
}
