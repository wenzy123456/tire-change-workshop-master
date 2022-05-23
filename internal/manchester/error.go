package manchester

import "fmt"

const (
	validationErrorCode      = "11"
	unAvailableTimeErrorCode = "22"
)

type tireChangeApplicationError struct {
	code  string
	error string
}

func (e tireChangeApplicationError) Error() string {
	return e.error
}

func newValidationError(cause error) *tireChangeApplicationError {
	return &tireChangeApplicationError{code: validationErrorCode, error: cause.Error()}
}

func newUnAvailableBookingError(e *tireChangeTimeEntity) *tireChangeApplicationError {
	return &tireChangeApplicationError{
		code:  unAvailableTimeErrorCode,
		error: fmt.Sprintf("tire change time %d is unavailable", e.ID)}
}
