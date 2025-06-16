package models

type (
	Response[T any] struct {
		Data T `json:"data"`
	}
)
