package auth

import (
	"backend/apperrors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
)

var TokenAuth *jwtauth.JWTAuth

func init() {
	godotenv.Load("../.env")
	secretKey := os.Getenv("SECRET_KEY")
	TokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)
}

func CreateToken(userId string, username string) (string, error) {
	claims := map[string]interface{}{"userID": userId, "username": username}

	jwtauth.SetExpiry(claims, time.Now().Add(15*time.Minute))
	_, tokenString, err := TokenAuth.Encode(claims)
	if err != nil {
		slog.Error("failed to create token", "error", err)
		return "", err
	}

	return tokenString, nil
}

type User struct {
	Username string
	ID       string
}

func ExtractToken(r *http.Request) (data *User, err error) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		slog.Debug("failed to decode token from request context", "error", err)
		return nil, err
	}

	data = &User{}
	if claims["username"] == nil || claims["userID"] == nil {
		return nil, apperrors.ErrInvalidToken
	}

	username, ok := claims["username"].(string)
	if !ok {
		slog.Warn("token username claim has an invalid type")
		return nil, apperrors.ErrTypeAssertion
	}

	id, ok := claims["userID"].(string)
	if !ok {
		slog.Warn("token user ID claim has an invalid type")
		return nil, apperrors.ErrTypeAssertion
	}

	data = &User{Username: username, ID: id}

	return data, nil
}
