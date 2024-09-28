package pkg

import (
	"crazyfarmbackend/src/constant"
	"errors"
	"fmt"
)

func PanicException_(key string, message string) {
	err := errors.New(message)
	err = fmt.Errorf("%s: %w", key, err)
	if err != nil {
		panic(err)
	}
}

func PanicException(responseKey constant.ResponseStatus, message string) {
	PanicException_(responseKey.GetResponseStatus(), message)
}
