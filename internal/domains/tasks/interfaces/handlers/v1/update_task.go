package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/supachai1998/task_services/internal/domains/tasks/models"
	"github.com/supachai1998/task_services/internal/entities"
	"github.com/supachai1998/task_services/internal/helpers"
	"gorm.io/gorm"
)

// GetTaskByID retrieves a task by its ID
// @Summary Retrieve a task by ID
// @Description Get a task using its unique ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.UpdateTaskRequest true "Task object"
// @Success 200 {object} models.ResponseSuccess{data=entities.Task} "Task found successfully"
// @Failure 400 {object} models.ResponseError "Invalid ID format"
// @Failure 404 {object} models.ResponseError "Task not found"
// @Router /v1/tasks/{id} [put]
func (h *Handler) UpdateTask(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.NewResponseError("Invalid ID format", "error"))
	}
	req := new(models.UpdateTaskRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.NewResponseError(err.Error(), "error"))
	}
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.NewResponseError(err.Error(), "error"))
	}

	var task entities.TaskUpdate
	copier.Copy(&task, req)

	task.Id = uint(id)

	if err := h.TaskUsecase.UpdateTask(&task); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, helpers.NewResponseError("Task not found", "error"))
		}
		return c.JSON(http.StatusInternalServerError, helpers.NewResponseError(err.Error(), "error"))
	}
	return c.JSON(http.StatusOK, helpers.NewResponseSuccess("Task updated", task))
}
