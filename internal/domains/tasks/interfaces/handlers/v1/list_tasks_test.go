package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	handlers "github.com/supachai1998/task_services/internal/domains/tasks/interfaces/handlers/v1"
	"github.com/supachai1998/task_services/internal/entities"
	"github.com/supachai1998/task_services/internal/interfaces"
	mocks "github.com/supachai1998/task_services/internal/mocks/tasks/usecases"
	"github.com/supachai1998/task_services/internal/models"
)

func TestListTasks(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockTaskUsecase(ctrl)
	handler := &handlers.Handler{TaskUsecase: mockUsecase}
	e := echo.New()
	e.Validator = interfaces.NewCustomValidator()

	t.Run("Success", func(t *testing.T) {
		// Define the expected tasks
		expectedTasks := []entities.Task{
			{
				Id:          1,
				Title:       "Task One",
				Description: "First task description",
				Status:      entities.TaskStatusInProgress,
			},
			{
				Id:          2,
				Title:       "Task Two",
				Description: "Second task description",
				Status:      entities.TaskStatusDone,
			},
		}

		// Expect the ListTasks method to be called and return the expected tasks without error
		mockUsecase.EXPECT().ListTasks().Return(expectedTasks, nil)

		// Create a new HTTP GET request
		req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks")

		// Invoke the handler
		if assert.NoError(t, handler.ListTasks(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusOK, rec.Code)

			// Unmarshal the response body
			var response models.ResponseSuccess
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Marshal and unmarshal to convert interface{} to []entities.Task
			tasksBytes, err := json.Marshal(response.Data)
			assert.NoError(t, err)

			var tasksResponse []entities.Task
			err = json.Unmarshal(tasksBytes, &tasksResponse)
			assert.NoError(t, err)

			// Assert the response message and data
			assert.Equal(t, "Tasks listed", response.Message)
			assert.Equal(t, expectedTasks, tasksResponse)
		}
	})

	t.Run("InternalServerError", func(t *testing.T) {
		// Define the error to be returned by the use case
		usecaseError := errors.New("database connection failed")

		// Expect the ListTasks method to be called and return an error
		mockUsecase.EXPECT().ListTasks().Return(nil, usecaseError)

		// Create a new HTTP GET request
		req := httptest.NewRequest(http.MethodGet, "/v1/tasks", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks")

		// Invoke the handler
		if assert.NoError(t, handler.ListTasks(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusInternalServerError, rec.Code)

			// Unmarshal the response body
			var response models.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, "database connection failed", response.Message)
		}
	})
}
