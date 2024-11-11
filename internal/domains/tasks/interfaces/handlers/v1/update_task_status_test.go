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
)

func TestUpdateTaskStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockTaskUsecase(ctrl)
	handler := &handlers.Handler{TaskUsecase: mockUsecase}
	e := echo.New()
	e.Validator = interfaces.NewCustomValidator()

	t.Run("Success", func(t *testing.T) {
		taskID := 1
		updateReq := models.UpdateTaskStatusRequest{
			Status: string(entities.TaskStatusDone),
		}

		// Expect the UpdateTaskStatus method to be called with the correct ID and status, returning no error
		mockUsecase.EXPECT().UpdateTaskStatus(uint(taskID), entities.TaskStatus(updateReq.Status)).Return(nil)

		// Marshal the request body to JSON
		reqBody, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		// Create a new HTTP PATCH request
		req := httptest.NewRequest(http.MethodPatch, "/v1/tasks/"+strconv.Itoa(taskID)+"/status", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id/status")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.UpdateTaskStatus(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusOK, rec.Code)

			// Unmarshal the response body
			var response globalModels.ResponseSuccess
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response message and data
			assert.Equal(t, "Task status updated", response.Message)
			assert.Equal(t, "", response.Data)
		}
	})

	t.Run("BadRequest_InvalidJSON", func(t *testing.T) {
		taskID := 2
		invalidJSON := `{"status": "done"` // Missing closing brace

		// Create a new HTTP PATCH request with invalid JSON
		req := httptest.NewRequest(http.MethodPatch, "/v1/tasks/"+strconv.Itoa(taskID)+"/status", bytes.NewReader([]byte(invalidJSON)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id/status")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.UpdateTaskStatus(c)) {
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

	t.Run("InternalServerError", func(t *testing.T) {
		taskID := 3
		updateReq := models.UpdateTaskStatusRequest{
			Status: string(entities.TaskStatusInProgress),
		}
		usecaseError := errors.New("database update failed")

		// Expect the UpdateTaskStatus method to be called with the correct ID and status, returning an error
		mockUsecase.EXPECT().UpdateTaskStatus(uint(taskID), entities.TaskStatus(updateReq.Status)).Return(usecaseError)

		// Marshal the request body to JSON
		reqBody, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		// Create a new HTTP PATCH request
		req := httptest.NewRequest(http.MethodPatch, "/v1/tasks/"+strconv.Itoa(taskID)+"/status", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id/status")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.UpdateTaskStatus(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusInternalServerError, rec.Code)

			// Unmarshal the response body
			var response globalModels.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, "database update failed", response.Message)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		invalidID := "abc"
		updateReq := models.UpdateTaskStatusRequest{
			Status: string(entities.TaskStatusDone),
		}

		// Marshal the request body to JSON
		reqBody, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		// Create a new HTTP PATCH request with an invalid ID
		req := httptest.NewRequest(http.MethodPatch, "/v1/tasks/"+invalidID+"/status", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id/status")
		c.SetParamNames("id")
		c.SetParamValues(invalidID)

		// Invoke the handler
		err = handler.UpdateTaskStatus(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		// Verify the response
		var response globalModels.ResponseError
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Invalid ID format", response.Message)

	})

	t.Run("BadRequest_ValidationFailure", func(t *testing.T) {
		taskID := 5
		updateReq := models.UpdateTaskStatusRequest{
			Status: "TO_DO",
		}

		// Marshal the request body to JSON
		reqBody, err := json.Marshal(updateReq)
		assert.NoError(t, err)

		// Create a new HTTP PATCH request
		req := httptest.NewRequest(http.MethodPatch, "/v1/tasks/"+strconv.Itoa(taskID)+"/status", bytes.NewReader(reqBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id/status")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.UpdateTaskStatus(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusBadRequest, rec.Code)

			// Unmarshal the response body
			var response globalModels.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Assert the response status and message
			assert.Equal(t, "error", response.Status)
			assert.Contains(t, response.Message, "invalid input")
		}
	})
}
