package status_errors

import (
	"errors"
	"fmt"
	terr "github.com/pkg/errors"
	"io"
)

var msgMap = map[ErrorType]string {
	UndefinedErr: "undefined error",
}

const (
	UndefinedErr = ErrorType(-1)
)

type ErrorType int

type httpError struct {
	errorType ErrorType
	cause error
	msg   string
}


func (e *httpError) Error() string { return e.msg + ": " + e.cause.Error() }
func (e *httpError) Cause() error  { return e.cause }

func (e httpError) Unwrap() error { return e.cause }

func (e *httpError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", e.Cause())
			io.WriteString(s, e.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, e.Error())
	}
}


func New(message string) error {
	err := &httpError{
		errorType: UndefinedErr,
		cause: errors.New(message),
		msg:   message,
	}
	return terr.WithStack(err)
}

func Wrap(err error, message string) error {
	return Wrapf(err, message)
}

func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	retErr := &httpError{
		errorType: UndefinedErr,
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	err =  Unwrap(err)
	if customErr, ok := err.(*httpError); ok {
		retErr.errorType = customErr.errorType
	}

	return terr.WithStack(retErr)
}

func Cause(err error) error {
	return terr.Cause(err)
}

func Unwrap(err error) error {
	type unwraper interface {
		Unwrap() error
	}

	for err != nil {
		unWrap, ok := err.(unwraper)
		if !ok {
			break
		}
		return unWrap.Unwrap()
	}

	return err
}

// NewTypeErr returns a err with the given error type code
// of course it contains a stack trace too
func NewTypeErr(message string, errorType ErrorType) error {
	err := &httpError{
		errorType: errorType,
		cause: errors.New(message),
		msg:   message,
	}
	return terr.WithStack(err)
}

func WrapTypeErr(err error, errorType ErrorType, message string) error {
	return WrapTypeErrf(err, errorType, message)
}

// WrapTypeErrf returns a err with the given error type code,
// no matter what error code it has brought before.
// of course it contains a stack trace too
func WrapTypeErrf(err error, errorType ErrorType, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	err = &httpError{
		errorType: errorType,
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	return terr.WithStack(err)
}


