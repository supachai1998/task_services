package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/supachai1998/task_services/internal/domains/tasks/usecases"
)

type Handler struct {
	TaskUsecase usecases.TaskUsecase
}

func NewTaskHandler(e *echo.Echo, taskUsecase usecases.TaskUsecase) {
	handler := &Handler{
		TaskUsecase: taskUsecase,
	}
	e.POST("/v1/tasks", handler.CreateTask)
	e.GET("/v1/tasks/:id", handler.GetTaskByID)
	e.PUT("/v1/tasks/:id", handler.UpdateTask)
	e.PATCH("/v1/tasks/:id/status", handler.UpdateTaskStatus)
	e.DELETE("/v1/tasks/:id", handler.DeleteTaskByID)
	e.GET("/v1/tasks", handler.ListTasks)
}
