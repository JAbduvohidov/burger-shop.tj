package errors

import "fmt"

type mError struct {
	text string
	err error
}

func ApiError(text string, err error) *mError {
	return &mError{text: text, err: err}
}

func (receiver *mError) Error() string {
	return fmt.Sprintf("error: %v", receiver.err.Error())
}

func (receiver *mError) Unwrap() error {
	return receiver.err
}
