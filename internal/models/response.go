package models

type ResponseSuccess struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Data    any    `json:"data"`
}

type ResponseError struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}
