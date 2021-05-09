package utils

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//SignToken signs a token using the secret from the config file
func SignToken(userID *string) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": *userID,
		"exp":    time.Now().Add(time.Second * 30 * 24 * 60 * 60),
	})

	tokenString, err := token.SignedString([]byte(GetConfig().Secret))

	return &tokenString, err
}

//ValidateToken validates that a given JWT token is valid
func ValidateToken(token *string) (*string, error) {
	jwtToken, err := jwt.Parse(*token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(GetConfig().Secret), nil
	})
	if err != nil {
		return nil, err
	}

	var userID string
	if claims, ok := jwtToken.Claims.(jwt.MapClaims); ok && jwtToken.Valid {
		expiredAt := claims["exp"].(string)
		now := time.Now()
		exp, _ := time.Parse(time.RFC3339, expiredAt)

		if !now.Before(exp) {
			return nil, fmt.Errorf("token expired")
		}

		userID = claims["userID"].(string)
	} else {
		return nil, fmt.Errorf("invalid claims")
	}

	return &userID, nil
}
