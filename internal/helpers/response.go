package helpers

import "github.com/supachai1998/task_services/internal/models"

func NewResponseSuccess(message string, data interface{}) models.ResponseSuccess {
	status := "success"
	return models.ResponseSuccess{
		Message: message,
		Status:  status,
		Data:    data,
	}
}

func NewResponseError(message string, status string) models.ResponseError {
	return models.ResponseError{
		Message: message,
		Status:  status,
	}
}
