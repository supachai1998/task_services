package handlers

import (
	"net/http"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/supachai1998/task_services/internal/domains/tasks/models"
	"github.com/supachai1998/task_services/internal/entities"
	"github.com/supachai1998/task_services/internal/helpers"
)

// CreateTask handles task creation
// @Summary Create a new task
// @Description Create a new task with the provided details
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.CreateTaskRequest true "Task object"
// @Success 201 {object} models.ResponseSuccess{data=entities.Task} "Task created successfully"
// @Failure 400 {object} models.ResponseError "Invalid input"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/tasks [post]
func (h *Handler) CreateTask(c echo.Context) error {
	req := new(models.CreateTaskRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.NewResponseError(err.Error(), "error"))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.NewResponseError(err.Error(), "error"))
	}

	task := new(entities.Task)
	copier.Copy(&task, req)
	if err := h.TaskUsecase.CreateTask(task); err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.NewResponseError(err.Error(), "error"))
	}
	return c.JSON(http.StatusCreated, helpers.NewResponseSuccess("Task created", task))
}
