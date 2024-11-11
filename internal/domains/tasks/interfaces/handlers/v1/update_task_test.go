package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	handlers "github.com/supachai1998/task_services/internal/domains/tasks/interfaces/handlers/v1"
	"github.com/supachai1998/task_services/internal/domains/tasks/models"
	"github.com/supachai1998/task_services/internal/entities"
	"github.com/supachai1998/task_services/internal/interfaces"
	mocks "github.com/supachai1998/task_services/internal/mocks/tasks/usecases"
	globalModels "github.com/supachai1998/task_services/internal/models"
	"gorm.io/gorm"
)

// Mocking copier.Copy is not straightforward since it's a function.
// To simulate copier.Copy failure, we would need to abstract it behind an interface.
// For simplicity, we'll skip the copier.Copy failure test as it requires refactoring the handler.

func TestUpdateTask(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockTaskUsecase(ctrl)
	handler := &handlers.Handler{TaskUsecase: mockUsecase}
	e := echo.New()
	e.Validator = interfaces.NewCustomValidator()

	t.Run("Success", func(t *testing.T) {
		taskID := 1
		updateReq := models.UpdateTaskRequest{
			Title:       "Updated Task Title",
			Description: "Updated Description",
		}

		// Expected task after update
		expectedTask := &entities.TaskUpdate{
			Id:          uint(taskID),
			Title:       &updateReq.Title,
			Description: &updateReq.Description,
		}

		// Set expectation: UpdateTask should be called with a TaskUpdate matching expectedTask
		mockUsecase.EXPECT().UpdateTask(gomock.AssignableToTypeOf(&entities.TaskUpdate{})).DoAndReturn(
			func(task *entities.TaskUpdate) error {
				// Verify that the task fields match the request
				assert.Equal(t, expectedTask.Id, task.Id)
				assert.Equal(t, expectedTask.Title, task.Title)
				assert.Equal(t, expectedTask.Description, task.Description)
				return nil
			},
		)

		// Marshal the request body to JSON
		reqBody, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		// Create a new HTTP PUT request
		req := httptest.NewRequest(http.MethodPut, "/v1/tasks/"+strconv.Itoa(taskID), bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.UpdateTask(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusOK, rec.Code)

			// Unmarshal the response body into models.ResponseSuccess
			var response globalModels.ResponseSuccess
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			taskBytes, err := json.Marshal(response.Data)
			assert.NoError(t, err)

			var taskResponse entities.Task
			err = json.Unmarshal(taskBytes, &taskResponse)
			assert.NoError(t, err)

			// Assert the response message and data
			assert.Equal(t, "Task updated", response.Message)
			assert.Equal(t, updateReq.Title, taskResponse.Title)
			assert.Equal(t, updateReq.Description, taskResponse.Description)
		}
	})
	t.Run("BadRequest_InvalidJSON", func(t *testing.T) {
		taskID := 2
		invalidJSON := `{"title": "Incomplete JSON"`

		// Create a new HTTP PUT request with invalid JSON
		req := httptest.NewRequest(http.MethodPut, "/v1/tasks/"+strconv.Itoa(taskID), bytes.NewReader([]byte(invalidJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.UpdateTask(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			// Unmarshal the response body
			var response globalModels.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Contains(t, response.Message, "unexpected")
		}
	})

	t.Run("InternalServerError_UsecaseFailure", func(t *testing.T) {
		taskID := 3
		updateReq := models.UpdateTaskRequest{
			Title:       "Another Update",
			Description: "Another Description",
		}

		// Expected task after update
		expectedTask := &entities.TaskUpdate{
			Id:          uint(taskID),
			Title:       &updateReq.Title,
			Description: &updateReq.Description,
		}

		// Define the use case error
		usecaseError := errors.New("database update failed")

		// Expect the UpdateTask method to be called with the correct task and return an error
		mockUsecase.EXPECT().UpdateTask(expectedTask).Return(usecaseError)

		// Marshal the request body to JSON
		reqBody, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		// Create a new HTTP PUT request
		req := httptest.NewRequest(http.MethodPut, "/v1/tasks/"+strconv.Itoa(taskID), bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.UpdateTask(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusInternalServerError, rec.Code)

			// Unmarshal the response body
			var response globalModels.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, usecaseError.Error(), response.Message)
		}
	})

	t.Run("InvalidID_NonInteger", func(t *testing.T) {
		invalidID := "abc"
		updateReq := models.UpdateTaskRequest{
			Title:       "Invalid ID Task",
			Description: "Description",
		}

		// Marshal the request body to JSON
		reqBody, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		// Create a new HTTP PUT request with an invalid ID
		req := httptest.NewRequest(http.MethodPut, "/v1/tasks/"+invalidID, bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(invalidID)

		// Invoke the handler
		if assert.NoError(t, handler.UpdateTask(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			// Unmarshal the response body
			var response globalModels.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, "Invalid ID format", response.Message)
		}
	})

	t.Run("NotFound_TaskDoesNotExist", func(t *testing.T) {
		taskID := 4
		updateReq := models.UpdateTaskRequest{
			Title:       "Non-Existent Task",
			Description: "Description",
		}

		// Expected task after update
		expectedTask := &entities.TaskUpdate{
			Id:          uint(taskID),
			Title:       &updateReq.Title,
			Description: &updateReq.Description,
		}

		// Define the use case error as a not found error
		usecaseError := gorm.ErrRecordNotFound

		// Expect the UpdateTask method to be called with the correct task and return a not found error
		mockUsecase.EXPECT().UpdateTask(expectedTask).Return(usecaseError)

		// Marshal the request body to JSON
		reqBody, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		// Create a new HTTP PUT request
		req := httptest.NewRequest(http.MethodPut, "/v1/tasks/"+strconv.Itoa(taskID), bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.UpdateTask(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusNotFound, rec.Code)

			// Unmarshal the response body
			var response globalModels.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, "Task not found", response.Message)
		}
	})

	t.Run("BadRequest_ValidationFailure", func(t *testing.T) {
		taskID := 5
		updateReq := models.UpdateTaskRequest{
			Title:       "Sh",
			Description: "Description",
		}

		// Marshal the request body to JSON
		reqBody, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		// Create a new HTTP PUT request
		req := httptest.NewRequest(http.MethodPut, "/v1/tasks/"+strconv.Itoa(taskID), bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.UpdateTask(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			// Unmarshal the response body
			var response globalModels.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Contains(t, response.Message, "title")
		}
	})
}
