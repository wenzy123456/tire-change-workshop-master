package london

import (
	"fmt"
	"time"
)

type validationError struct {
	error
}

type unAvailableBookingError struct {
	error string
}

func newUnAvailableBookingError(e *tireChangeTimeEntity) unAvailableBookingError {
	return unAvailableBookingError{error: fmt.Sprintf("tire change time %s is unavailable", e.UUID)}
}

func (e unAvailableBookingError) Error() string {
	return e.error
}

type invalidTireChangeTimesPeriodError struct {
	error string
}

func newInvalidTirChangeTimesPeriodError(from time.Time, until time.Time) invalidTireChangeTimesPeriodError {
	return invalidTireChangeTimesPeriodError{
		error: fmt.Sprintf("cannot fetch tire change times with invalid date period: %s - %s", from, until),
	}
}

func (e invalidTireChangeTimesPeriodError) Error() string {
	return e.error
}
