package controller

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type UserClaims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func generateToken(
	userId int, username, role string,
	secretKey []byte,
	expier time.Duration,
) (string, error) {
	claims := UserClaims{
		ID:       userId,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expier)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateUser(c echo.Context) (*UserClaims, error) {
	userInterface := c.Get("user")
	if userInterface == nil {
		return nil, fmt.Errorf("user not found")
	}

	user, ok := userInterface.(*jwt.Token)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok || !user.Valid {
		fmt.Printf("%T\n", user.Claims)
		return nil, fmt.Errorf("user not found")
	}

	userId, ok := claims["id"].(float64)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return &UserClaims{
		ID:       int(userId),
		Username: username,
		Role:     role,
	}, nil
}
