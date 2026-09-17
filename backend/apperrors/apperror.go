package apperrors

import "errors"

var (
	ErrUserNotFound         = errors.New("user_not_found")
	ErrInvalidCredentials   = errors.New("invalid_credentials")
	ErrUserAlreadyExists    = errors.New("user_already_exists")
	ErrMobileAlreadyExists  = errors.New("mobile_already_exists")
	ErrNotAFriend           = errors.New("unauthorized attempt - sending msg to an unknown person")
	ErrInvalidMessage       = errors.New("invalid_message")
	ErrInvalidFriendRequest = errors.New("invalid_friend_request")
	ErrFriendRequestExists  = errors.New("friend_request_exists")

	ErrTypeAssertion = errors.New("assertion_failed")

	ErrInvalidToken = errors.New("invalid_token")

	ErrEmptyRedisClient = errors.New("redis not started")
	ErrEmptySqlClient   = errors.New("sql not started")
)
