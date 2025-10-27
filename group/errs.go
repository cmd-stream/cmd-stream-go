package group

import "errors"

const errorPrefix = "cmdstream group: "

func WrapError(err error) error {
	return errors.New(errorPrefix + err.Error())
}
