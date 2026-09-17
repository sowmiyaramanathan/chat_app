package auth

import (
	"backend/apperrors"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

var TokenAuth *jwtauth.JWTAuth

func Init() error {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return errors.New("SECRET_KEY must be set")
	}

	TokenAuth = jwtauth.New("HS256", []byte(secretKey), nil)
	return nil
}

func CreateToken(userId string, username string) (string, string, error) {
	if TokenAuth == nil {
		return "", "", errors.New("token auth is not initialized")
	}
	accessClaims := map[string]interface{}{
		"userID":    userId,
		"username":  username,
		"tokenType": "access",
		"jti":       uuid.NewString(),
	}

	jwtauth.SetExpiry(accessClaims, time.Now().Add(15*time.Minute))
	_, accessToken, err := TokenAuth.Encode(accessClaims)
	if err != nil {
		slog.Error("failed to create acess token", "error", err)
		return "", "", err
	}

	refreshClaims := map[string]interface{}{
		"userID":    userId,
		"username":  username,
		"tokenType": "refresh",
		"jti":       uuid.NewString(),
	}
	jwtauth.SetExpiry(refreshClaims, time.Now().Add(7*24*time.Hour))
	_, refreshToken, err := TokenAuth.Encode(refreshClaims)
	if err != nil {
		slog.Error("failed to create refresh token", "error", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateRefreshToken verifies that token is a valid, unexpired JWT issued for
// refreshes. Database session validation is intentionally performed by the
// caller so that refresh-token rotation can invalidate previous tokens.
func ValidateRefreshToken(tokenString string) (*User, error) {
	if TokenAuth == nil {
		return nil, errors.New("token auth is not initialized")
	}

	token, err := jwtauth.VerifyToken(TokenAuth, tokenString)
	if err != nil || token == nil {
		return nil, apperrors.ErrInvalidToken
	}

	tokenType, ok := token.Get("tokenType")
	if !ok || tokenType != "refresh" {
		return nil, apperrors.ErrInvalidToken
	}

	username, ok := token.Get("username")
	if !ok {
		return nil, apperrors.ErrInvalidToken
	}
	usernameString, ok := username.(string)
	if !ok {
		return nil, apperrors.ErrTypeAssertion
	}

	userID, ok := token.Get("userID")
	if !ok {
		return nil, apperrors.ErrInvalidToken
	}
	userIDString, ok := userID.(string)
	if !ok {
		return nil, apperrors.ErrTypeAssertion
	}

	return &User{Username: usernameString, ID: userIDString}, nil
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
	if claims["tokenType"] != "access" {
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
