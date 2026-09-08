package auth

import (
	"backend/apperrors"
	"errors"
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
	if secretKey == "" {
		slog.Error("SECRET_KEY is empty")
	}

	TokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)
}

func CreateToken(userId string, username string) (string, string, error) {
	if TokenAuth == nil {
		return "", "", errors.New("token auth is not initialized")
	}
	claims := map[string]interface{}{"userID": userId, "username": username}

	jwtauth.SetExpiry(claims, time.Now().Add(15*time.Minute))
	_, accessToken, err := TokenAuth.Encode(claims)
	if err != nil {
		slog.Error("failed to create acess token", "error", err)
		return "", "", err
	}

	jwtauth.SetExpiry(claims, time.Now().Add(7*24*time.Hour))
	_, refreshToken, err := TokenAuth.Encode(claims)
	if err != nil {
		slog.Error("failed to create refresh token", "error", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

type User struct {
	Username string
	ID       string
}

func ExtractToken(r *http.Request) (data *User, err error) {
	if TokenAuth == nil {
		return nil, errors.New("token auth is not initialized")
	}
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
