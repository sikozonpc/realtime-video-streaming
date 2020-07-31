package user

// User data
type User struct {
	ID       int
	Username string `json:"username"`
	Email    string `json:"email"`
}

// type UserService interface {
// 	User(id int) (*User, error)
// 	Users() ([]*User, error)
// 	CreateUser(u *User) error
// 	DeleteUser(id int) error
// }
