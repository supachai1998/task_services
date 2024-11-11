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
	"github.com/supachai1998/task_services/internal/interfaces"
	mocks "github.com/supachai1998/task_services/internal/mocks/tasks/usecases"
	"github.com/supachai1998/task_services/internal/models"
)

func TestDeleteTaskByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mocks.NewMockTaskUsecase(ctrl)
	handler := &handlers.Handler{TaskUsecase: mockUsecase}
	e := echo.New()
	e.Validator = interfaces.NewCustomValidator()

	t.Run("Success", func(t *testing.T) {
		taskID := 1

		// Expect the DeleteTaskByID method to be called with the correct ID and return no error
		mockUsecase.EXPECT().DeleteTaskByID(uint(taskID)).Return(nil)

		// Create a new HTTP DELETE request
		req := httptest.NewRequest(http.MethodDelete, "/v1/tasks/"+strconv.Itoa(taskID), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		// Invoke the handler
		if assert.NoError(t, handler.DeleteTaskByID(c)) {
			// Assert the HTTP status code
			assert.Equal(t, http.StatusNoContent, rec.Code)

			// Verify the response body is empty
			assert.Empty(t, rec.Body.String())
		}
	})

	t.Run("TaskNotFound", func(t *testing.T) {
		taskID := 2
		mockUsecase.EXPECT().DeleteTaskByID(uint(taskID)).Return(errors.New("task not found"))

		req := httptest.NewRequest(http.MethodDelete, "/v1/tasks/"+strconv.Itoa(taskID), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/v1/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.Itoa(taskID))

		if assert.NoError(t, handler.DeleteTaskByID(c)) {
			assert.Equal(t, http.StatusNotFound, rec.Code)

			var response models.ResponseError
			err := json.Unmarshal(rec.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, "error", response.Status)
			assert.Equal(t, "Task not found", response.Message)
		}
	})

	t.Run("InvalidID", func(t *testing.T) {
		// Prepare the request with an invalid ID
		req := httptest.NewRequest(http.MethodDelete, "/tasks/invalid", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/tasks/:id")
		c.SetParamNames("id")
		c.SetParamValues("invalid")

		// Execute the handler
		err := handler.DeleteTaskByID(c)
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
