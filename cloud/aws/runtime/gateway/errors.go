package gateway

import "fmt"

// An error indicating the JSON event type is not handled by this lambda gateway
type UnhandledLambdaEventError struct {
	Message string
	Cause   error
}

func (e *UnhandledLambdaEventError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s", e.Message, e.Cause.Error())
	}
	return e.Message
}

func (e *UnhandledLambdaEventError) Unwrap() error {
	return e.Cause
}

func NewUnhandledLambdaEventError(cause error) *UnhandledLambdaEventError {
	return &UnhandledLambdaEventError{
		Message: "the nitric lambda gateway does not handle this lambda event type",
		Cause:   cause,
	}
}
