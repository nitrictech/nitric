package errors

import "fmt"

type FatalError struct {
	msg string
}

func (fe *FatalError) Error() string {
	return fmt.Sprintf("fatal: %s", fe.msg)
}

var _ error = &FatalError{}

func Fatal(msg string) *FatalError {
	return &FatalError{msg: msg}
}
