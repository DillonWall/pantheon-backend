package auth

import (
	"fmt"
	"pantheon-auth/graph/model"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository struct {
	// todo: should use a database and instead do a query on it instead of an array in a real-world environment
	Users []*model.User
}

func (ur *UserRepository) GetUser_byUsername(username string) (*model.User, error) {
	var user *model.User
	for _, u := range ur.Users {
		if u.Username == username {
			user = u
			break
		}
	}

	if user == nil {
		return nil, fmt.Errorf("User not found")
	}

	return user, nil
}

func (ur *UserRepository) CreateUser(username string, password string) (error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Failed to hash the password: %w", err)
	}

	user := &model.User{
		Username:     username,
		Passwordhash: string(hashedPassword),
	}
	ur.Users = append(ur.Users, user)

	return nil
}
