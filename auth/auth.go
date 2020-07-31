package auth

import (
	"goproject/user"
	"log"
)

// Register a new user
func (a Auth) Register(creds Credentials) error {

	// Save to the database
	sqlStatment := `
	INSERT INTO users(username, password, email)
	VALUES ($1, $2, $3)
	`
	res, err := a.DB.Exec(sqlStatment, creds.Username, creds.Password, creds.Email)
	if err != nil {
		return err
	}

	log.Print("Hello there")
	log.Print(res)

	return nil
}

// GetUsers gets all users
func (a Auth) GetUsers() ([]user.User, error) {

	var users []user.User
	var email string
	var id string
	var username string

	sqlStatment := `SELECT id, email, username FROM users`
	a.DB.QueryRow(sqlStatment).Scan(&id, &email, &username)

	return users, nil
}
