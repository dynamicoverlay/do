package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"imabad.dev/do/lib/utils"
)

type ContextKey string

type UnauthenticatedError struct {
	message string `json:"message"`
	code    int    `json:"code"`
}

func Unautherror(w http.ResponseWriter) {
	jsonBytes, err := json.Marshal(UnauthenticatedError{
		message: "Unauthenticated",
		code:    401,
	})
	if err != nil {
		w.Write(jsonBytes)
	} else {
		w.Write(make([]byte, 0))
	}
}

func AuthCheckMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userID, err := validateToken(ctx, r)
		if err != nil {
			Unautherror(w)
			return
		}

		if userID != nil {
			ctx = context.WithValue(ctx, ContextKey("UserID"), userID)
		}
		fmt.Println("User with id ", userID)

		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateToken(ctx context.Context, r *http.Request) (*int, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, nil
	}

	userIDString, err := utils.ValidateToken(&tokenString)
	userID, err := strconv.Atoi(utils.StrValue(userIDString))
	return &userID, err
}
