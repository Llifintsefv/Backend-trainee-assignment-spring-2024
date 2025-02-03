package model

import "errors"

type ErrorResponse struct {
	Reason string `json:"reason"`
}

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrForbidden       = errors.New("forbidden")
	ErrTenderNotFound  = errors.New("tender not found")
	ErrBidNotFound     = errors.New("bid not found")
	ErrVersionNotFound = errors.New("version not found")
	ErrDecisionSubmit  = errors.New("decision cannot be submitted")
	ErrFeedbackSubmit  = errors.New("feedback cannot be submitted")
)
