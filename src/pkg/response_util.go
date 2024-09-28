package pkg

import (
	"crazyfarmbackend/src/constant"
	"crazyfarmbackend/src/domain/dto"
)

func Null() interface{} {
	return nil
}

func BuildResponse[T any](responseStatus constant.ResponseStatus, data T) dto.ApiResponse[T] {
	return BuildResponse_(responseStatus.GetResponseStatus(), data)
}

func BuildResponse_[T any](status string, data T) dto.ApiResponse[T] {
	return dto.ApiResponse[T]{
		ResponseKey: status,
		Data:        data,
	}
}
