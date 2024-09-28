package dto

type ApiResponse[T any] struct {
	ResponseKey string `json:"response_key"`
	Data        T      `json:"data"`
}
