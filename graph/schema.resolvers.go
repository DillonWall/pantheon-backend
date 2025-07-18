package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.76

import (
	"context"
	"fmt"
	"log"
	"os"
	"pantheon-auth/graph/model"
	"pantheon-auth/pkg/auth"
	"pantheon-auth/pkg/imageapi"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, input model.UserData) (*model.AuthResponse, error) {
	// Validation
    if input.Username == "" || input.Password == "" {
        return nil, fmt.Errorf("Username and password are required")
    }
	user, err := r.UserRepo.GetUser_byUsername(input.Username)
	if user != nil {
		return nil, fmt.Errorf("Username already taken")
	}

	err = r.UserRepo.CreateUser(input.Username, input.Password)
	if err != nil {
		return nil, err
	}

	tokenString, err := auth.GenerateToken(input.Username)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: &tokenString,
	}, nil
}

// Login is the resolver for the login field.
func (r *mutationResolver) Login(ctx context.Context, input model.UserData) (*model.AuthResponse, error) {
	// Validation
    if input.Username == "" || input.Password == "" {
        return nil, fmt.Errorf("Username and password are required")
    }
	user, err := r.UserRepo.GetUser_byUsername(input.Username)
	if err != nil {
		return nil, err
	}

	// Hash the password from the request and compare it with the stored password
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Passwordhash),
		[]byte(input.Password),
	)
	if err != nil {
		return nil, fmt.Errorf("Invalid username or password: %w", err)
	}

	tokenString, err := auth.GenerateToken(input.Username)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		Token: &tokenString,
	}, nil
}

// Verify is the resolver for the verify field.
func (r *mutationResolver) Verify(ctx context.Context, token string) (bool, error) {
	return auth.ValidateToken(token, r.UserRepo)
}

// SearchImages is the resolver for the searchImages field.
func (r *queryResolver) SearchImages(ctx context.Context, token string, query string) ([]*model.Image, error) {
	valid, err := auth.ValidateToken(token, r.UserRepo)
	if !valid || err != nil {
		return nil, err
	}

	timeoutSec, err := strconv.Atoi(os.Getenv("IMAGE_API_TIMEOUT_SEC"))
	if err != nil {
		return nil, fmt.Errorf("Server environment variable error: %w", err)
	}

	// Concurrent image fetching
	type result struct {
		image *model.Image
		err   error
	}

	results := make(chan result)
	ctx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSec)*time.Second)
	defer cancel()

	// Launch goroutines for each API
	for _, api := range r.ImageAPIs {
		go func(a imageapi.API) {
			image, err := a.SearchSingleImage(ctx, query)
			results <- result{image, err}
		}(api)
	}

	// Collect results
	var allImages []*model.Image
	numAPIs := len(r.ImageAPIs)
	for i := 0; i < numAPIs; i++ {
		res := <-results
		if res.err != nil {
			log.Printf("Error querying %T: %v", res.err, res.err)
			continue
		}
		allImages = append(allImages, res.image)
	}
	close(results)

	if len(allImages) == 0 {
		return nil, fmt.Errorf("no images found for query: %s", query)
	}

	return allImages, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
