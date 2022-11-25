package utils

import (

	"github.com/pkg/errors"
)

//ErrorString ...
func ErrorString(e error, s string) string{
	err := errors.Wrap(e, s)

	return err.Error()
}

//ErrorNew ...
func ErrorNew(s string) error{
	err := errors.New(s)

	return err
}

//ErrorWrap ...
func ErrorWrap(e error, s string) error{
	err := errors.Wrap(e, s)

	return err
}