package apperrors

import "errors"

var (
	ErrUserNotFound        = errors.New("user_not_found")
	ErrInvalidCredentials  = errors.New("invalid_credentials")
	ErrUserAlreadyExists   = errors.New("user_already_exists")
	ErrMobileAlreadyExists = errors.New("mobile_already_exists")

	ErrTypeAssertion = errors.New("assertion_failed")

	ErrInvalidToken = errors.New("invalid_token")
)
