package auth

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateToken(username string) (string, error) {
	// Generate a JWT token
	tokenExpireTimeHours, err := strconv.Atoi(os.Getenv("TOKEN_EXPIRE_TIME_HOURS"))
	if err != nil {
		return "", fmt.Errorf("Server environment variable error: %w", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(tokenExpireTimeHours)).Unix(),
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("Failed to generate JWT token: %w", err)
	}

	return tokenString, nil
}

func ValidateToken(token string, userRepo *UserRepository) (bool, error) {
	// Parse the JWT token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if err != nil {
		return false, fmt.Errorf("invalid token: %w", err)
	}

	if !parsedToken.Valid {
		return false, nil
	}

	// Extract claims from the token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, fmt.Errorf("invalid token claims")
	}

	// Check if the username exists in the claims
	username, ok := claims["username"].(string)
	if !ok || username == "" {
		return false, fmt.Errorf("invalid username claim")
	}

	// Verify the user exists
    user, err := userRepo.GetUser_byUsername(username)
    if err != nil {
        return false, err
    }

	if user == nil {
		return false, fmt.Errorf("user not found for the token")
	}

	return true, nil
}
