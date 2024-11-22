package remi_errors

import "errors"

var NotADurationError = errors.New("Not a duration")
var InvalidTimestampError = errors.New("Invalid timestamp format")
var MissingMessageError = errors.New("Missing message")
var InvalidReminderIdError = errors.New("Invalid reminder ID")
var MissingReminderIdError = errors.New("Missing reminder ID")
