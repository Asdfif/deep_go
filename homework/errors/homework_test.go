package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	if len(e.errs) == 0 {
		return ""
	}

	msg := fmt.Sprintf("%d errors occured:\n", len(e.errs))

	for _, err := range e.errs {
		msg += fmt.Sprintf("\t* %v", err)
	}

	msg += "\n"

	return msg
}

func Append(err error, errs ...error) *MultiError {
	if err == nil {
		return Append(&MultiError{}, errs...)
	}

	multiError, _ := err.(*MultiError)
	multiError.errs = append(multiError.errs, errs...)
	return multiError
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occured:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
