package user

// User data
type User struct {
	ID       int
	Username string `json:"username"`
	Email    string `json:"email"`
}
