package app

import (
	"context"
	"errors"
	"net/http"
	"os"

	"github.com/Meexe/videocall/models"
	u "github.com/Meexe/videocall/utils"
	"github.com/dgrijalva/jwt-go"
)

func ValidateToken(str string) (*models.Token, error) {

	tk := &models.Token{}

	token, err := jwt.ParseWithClaims(str, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	if err != nil {
		err := errors.New("Malformed authentication token")
		return nil, err
	}

	if !token.Valid {
		err := errors.New("Token is not valid")
		return nil, err
	}

	return tk, err
}

var Authentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/user/new", "/api/user/login"}
		requestPath := r.URL.Path

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenCookie, err := r.Cookie("jwt")

		if err != nil {
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		tk, err := ValidateToken(tokenCookie.Value)

		if err != nil {
			response = u.Message(false, err.Error())
			w.WriteHeader(http.StatusForbidden)
			u.Respond(w, response)
			return
		}

		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
