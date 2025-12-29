package utils

import "fmt"

func WrapErr(format string, wrappedErr, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%w: "+format+": %v", wrappedErr, err)
}

func NewErr(format string, wrappedErr error, args ...interface{}) error {
	msg := fmt.Sprintf(format, args...)
	return fmt.Errorf("%w: %s", wrappedErr, msg)
}
