package error_factory

import (
	"errors"
	"fmt"
)

func NewError(cause string) error {
	return errors.New(cause)
}

func NewErrorWithMessage(cause string, message any) error {
	return fmt.Errorf("%s: %v", cause, message)
}
