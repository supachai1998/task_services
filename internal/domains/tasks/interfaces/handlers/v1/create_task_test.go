package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/supachai1998/task_services/internal/domains/tasks/interfaces/handlers/v1"
	taskModels "github.com/supachai1998/task_services/internal/domains/tasks/models"
	"github.com/supachai1998/task_services/internal/entities"
	"github.com/supachai1998/task_services/internal/interfaces"
	mocks "github.com/supachai1998/task_services/internal/mocks/tasks/usecases"
	"github.com/supachai1998/task_services/internal/models"
)

func TestCreateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockTaskUsecase(ctrl)
	handler := &handlers.Handler{TaskUsecase: mockUsecase}
	e := echo.New()
	e.Validator = interfaces.NewCustomValidator()

	t.Run("Success", func(t *testing.T) {
		// Define the input request
		createTaskReq := &taskModels.CreateTaskRequest{
			Title:       "Test Task",
			Description: "Test Description",
		}
		payload, err := json.Marshal(createTaskReq)
		assert.NoError(t, err)

		// Define the expected task that should be passed to the use case
		expectedTask := &entities.Task{
			Title:       createTaskReq.Title,
			Description: createTaskReq.Description,
		}

		// Set expectation: CreateTask should be called with a task matching expectedTask
		mockUsecase.EXPECT().CreateTask(gomock.AssignableToTypeOf(&entities.Task{})).DoAndReturn(
			func(task *entities.Task) error {
				// Verify that the task fields match the request
				assert.Equal(t, expectedTask.Title, task.Title)
				assert.Equal(t, expectedTask.Description, task.Description)
				assert.Equal(t, expectedTask.Status, task.Status)
				return nil
			},
		)

		// Create a new HTTP POST request
		req := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBuffer(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Invoke the handler
		if assert.NoError(t, handler.CreateTask(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusCreated, rec.Code)

			// Unmarshal the response body into models.ResponseSuccess
			var response models.ResponseSuccess
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Marshal and unmarshal to convert interface{} to entities.Task
			taskBytes, err := json.Marshal(response.Data)
			assert.NoError(t, err)

			var taskResponse entities.Task
			err = json.Unmarshal(taskBytes, &taskResponse)
			assert.NoError(t, err)

			// Assert the response message and data
			assert.Equal(t, "Task created", response.Message)
			assert.Equal(t, expectedTask.Title, taskResponse.Title)
			assert.Equal(t, expectedTask.Description, taskResponse.Description)
			assert.Equal(t, expectedTask.Status, taskResponse.Status)
		}
	})

	t.Run("BindError", func(t *testing.T) {
		// Define an invalid JSON payload
		invalidPayload := []byte(`{"title": "Incomplete JSON`)

		// Create a new HTTP POST request with invalid JSON
		req := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBuffer(invalidPayload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Invoke the handler
		if assert.NoError(t, handler.CreateTask(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			// Unmarshal the response body into models.ResponseError
			var response models.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Contains(t, response.Message, "unexpected")
		}
	})

	t.Run("UsecaseError", func(t *testing.T) {
		// Define the input request
		createTaskReq := &taskModels.CreateTaskRequest{
			Title:       "Test Task",
			Description: "Test Description",
		}
		payload, err := json.Marshal(createTaskReq)
		assert.NoError(t, err)

		// Define the expected task that should be passed to the use case
		expectedTask := &entities.Task{
			Title:       createTaskReq.Title,
			Description: createTaskReq.Description,
		}

		// Define the error to be returned by the use case
		usecaseError := assert.AnError

		// Set expectation: CreateTask should be called with a task matching expectedTask and return an error
		mockUsecase.EXPECT().CreateTask(gomock.AssignableToTypeOf(&entities.Task{})).DoAndReturn(
			func(task *entities.Task) error {
				// Verify that the task fields match the request
				assert.Equal(t, expectedTask.Title, task.Title)
				assert.Equal(t, expectedTask.Description, task.Description)
				assert.Equal(t, expectedTask.Status, task.Status)
				return usecaseError
			},
		)

		// Create a new HTTP POST request
		req := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBuffer(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Invoke the handler
		if assert.NoError(t, handler.CreateTask(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusInternalServerError, rec.Code)

			// Unmarshal the response body into models.ResponseError
			var response models.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, usecaseError.Error(), response.Message)
		}
	})

	t.Run("ValidateError", func(t *testing.T) {
		// Define the input request with an empty title
		createTaskReq := &taskModels.CreateTaskRequest{
			Title:       "",
			Description: "Test Description",
		}
		payload, err := json.Marshal(createTaskReq)
		assert.NoError(t, err)

		// Create a new HTTP POST request
		req := httptest.NewRequest(http.MethodPost, "/v1/tasks", bytes.NewBuffer(payload))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// Invoke the handler
		if assert.NoError(t, handler.CreateTask(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			// Unmarshal the response body into models.ResponseError
			var response models.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Contains(t, response.Message, "title")
		}
	})
}
