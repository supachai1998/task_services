package handlers_test

import (
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
	"github.com/supachai1998/task_services/internal/entities"
	"github.com/supachai1998/task_services/internal/interfaces"
	mocks "github.com/supachai1998/task_services/internal/mocks/tasks/usecases"
	"github.com/supachai1998/task_services/internal/models"
)

func TestGetTaskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockTaskUsecase(ctrl)
	handler := &handlers.Handler{TaskUsecase: mockUsecase}
	e := echo.New()
	e.Validator = interfaces.NewCustomValidator()

	t.Run("Success", func(t *testing.T) {
		taskID := 1
		expectedTask := &entities.Task{
			Id:          uint(taskID),
			Title:       "Test Task",
			Description: "This is a test task",
			Status:      entities.TaskStatusDone,
		}

		// Expect the GetTaskByID to be called with the correct ID and return the expected task without error
		mockUsecase.EXPECT().GetTaskByID(uint(taskID)).Return(expectedTask, nil)

		// Create a new HTTP GET request
		req := httptest.NewRequest(http.MethodGet, "/v1/tasks/"+strconv.Itoa(taskID), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.GetTaskByID(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusOK, rec.Code)

			// Unmarshal the response body
			var response models.ResponseSuccess
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			taskBytes, err := json.Marshal(response.Data)
			assert.NoError(t, err)

			var taskResponse entities.Task
			err = json.Unmarshal(taskBytes, &taskResponse)
			assert.NoError(t, err)

			// Assert the response message and data
			assert.Equal(t, "Task found", response.Message)
			assert.Equal(t, expectedTask.Id, taskResponse.Id)
			assert.Equal(t, expectedTask.Title, taskResponse.Title)
			assert.Equal(t, expectedTask.Description, taskResponse.Description)
			assert.Equal(t, expectedTask.Status, taskResponse.Status)
		}
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		taskID := 2

		// Expect the GetTaskByID to be called with the correct ID and return an error
		mockUsecase.EXPECT().GetTaskByID(uint(taskID)).Return(nil, errors.New("task not found"))

		// Create a new HTTP GET request
		req := httptest.NewRequest(http.MethodGet, "/v1/tasks/"+strconv.Itoa(taskID), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.GetTaskByID(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusNotFound, rec.Code)

			// Unmarshal the response body
			var response models.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, "Task not found", response.Message)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		invalidID := "abc"

		// Create a new HTTP GET request with an invalid ID
		req := httptest.NewRequest(http.MethodGet, "/v1/tasks/"+invalidID, nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(invalidID)

		// Invoke the handler
		err := handler.GetTaskByID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Verify the response
		var response models.ResponseError
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Invalid ID format", response.Message)
	})
}
