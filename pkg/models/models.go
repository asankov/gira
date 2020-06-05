package models

// Game is the representation of a game
// in the database.
type Game struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status"`
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

// Status is the type that represents the status of a game
type Status string

var (
	// StatusTODO is the To Do status of a game
	StatusTODO Status = "To Do"
	// StatusInProgress is the In progress status of a game
	StatusInProgress Status = "In Progress"
	// StatusDone is the Done status of the game
	StatusDone Status = "Done"
)

// UserGame is the representation of a user game relation
// in the database.
type UserGame struct {
	ID     string `json:"id"`
	User   *User  `json:"user"`
	Game   *Game  `json:"game"`
	Status Status `json:"status"`
}

// UserResponse is the response that is returned
// when a user is logged in.
type UserResponse struct {
	Token string `json:"token"`
}

type UserGameResponse struct {
	Games []*Game `json:"games"`
}
