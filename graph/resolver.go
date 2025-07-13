package graph

import "pantheon-auth/pkg/auth"

//go:generate go run github.com/99designs/gqlgen generate

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.


type Resolver struct {
	// todo: add database instead of array
    UserRepo *auth.UserRepository
}
