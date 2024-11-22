package model

type CommonResponse[T any] struct {
	HttpStatus int    `json:"httpStatus"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}
