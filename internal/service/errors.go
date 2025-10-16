package service

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidUsername   = errors.New("invalid username")
	ErrEmptyUserID       = errors.New("user id cannot be empty")
)

var (
	ErrProgressNotFound      = errors.New("user progress not found")
	ErrProgressAlreadyExists = errors.New("user progress already exists")
	ErrInvalidProgressData   = errors.New("invalid progress data")
	ErrEmptyProgressID       = errors.New("progress id cannot be empty")
	ErrEmptyWordID           = errors.New("word id cannot be empty")
)

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrSessionAlreadyExists = errors.New("session already exists")
	ErrInvalidSessionStatus = errors.New("invalid session status")
	ErrEmptySessionID       = errors.New("session id cannot be empty")
)

var (
	ErrStatisticsNotFound    = errors.New("statistics not found")
	ErrInvalidStatisticsData = errors.New("invalid statistics data")
)

var (
	ErrWordSetNotFound      = errors.New("word set not found")
	ErrWordSetAlreadyExists = errors.New("word set already exists")
	ErrEmptyWordSetID       = errors.New("word set id cannot be empty")
)
