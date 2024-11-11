package handlers_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	handlers "github.com/supachai1998/task_services/internal/domains/tasks/interfaces/handlers/v1"
	mocks "github.com/supachai1998/task_services/internal/mocks/tasks/usecases"
)

func TestNewTaskHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	mockUsecase := mocks.NewMockTaskUsecase(ctrl)

	handlers.NewTaskHandler(e, mockUsecase)

	routes := e.Routes()

	expectedRoutes := []struct {
		Method string
		Path   string
	}{
		{"POST", "/v1/tasks"},
		{"GET", "/v1/tasks/:id"},
		{"PUT", "/v1/tasks/:id"},
		{"PATCH", "/v1/tasks/:id/status"},
		{"DELETE", "/v1/tasks/:id"},
		{"GET", "/v1/tasks"},
	}

	for _, er := range expectedRoutes {
		found := false
		for _, r := range routes {
			if r.Method == er.Method && r.Path == er.Path {
				found = true
				break
			}
		}
		assert.True(t, found, "Route not registered: %s %s", er.Method, er.Path)
	}
}
