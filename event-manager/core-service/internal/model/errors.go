package model

import "errors"

var (
	ErrUserNotFound             = errors.New("user not found")
	ErrUserAlreadyAnParticipant = errors.New("user is already a participant of event")
)
