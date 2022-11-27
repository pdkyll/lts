package utils

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Define the errors we want to use inside the business logic & domain.
var (
	ErrUnmatchPassword    	= errors.New("unmatch password ")
	ErrIncorrectPassword   	= errors.New("password is required (min: 6 characters)")
	ErrIncorrectEmail      	= errors.New("incorrect email")
	ErrIncorrectDomain     	= errors.New("incorrect domain")
	ErrUserExists      		= errors.New("user already exists")
	ErrUserNotFound    		= errors.New("user not found")
	ErrDBError            	= errors.New("data base error")
	ErrNotFound            	= errors.New("resource not found")
	ErrConflict            	= errors.New("datamodel conflict")
	ErrForbidden           	= errors.New("forbidden access")
	ErrReviewInput         	= errors.New("review your input")
	ErrBadToken          	= errors.New("bad token")
	ErrBadRequest          	= errors.New("bad request")
	ErrDuplicateEntry      	= errors.New("duplicate entry")
	ErrGone                	= errors.New("resource gone")
	ErrNotAllowed          	= errors.New("operation not allowed")
	ErrBusy                	= errors.New("resource is busy")
	ErrUnauthorized        	= errors.New("unauthorized")
)

//LogError ...
func LogError(log *zap.Logger, err error) {
	log.Error("test", 
		//zap.String("err", err.Error()),
		//zap.String("ctx", string(ctx.Err().Error())),
	)
}

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