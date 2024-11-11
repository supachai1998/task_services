package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/supachai1998/task_services/internal/helpers"
)

// GetTaskByID retrieves a task by its ID
// @Summary Retrieve a task by ID
// @Description Get a task using its unique ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.ResponseSuccess{data=entities.Task} "Task found successfully"
// @Failure 400 {object} models.ResponseError "Invalid ID format"
// @Failure 404 {object} models.ResponseError "Task not found"
// @Router /v1/tasks/{id} [get]
func (h *Handler) GetTaskByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.NewResponseError("Invalid ID format", "error"))
	}
	task, err := h.TaskUsecase.GetTaskByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, helpers.NewResponseError("Task not found", "error"))
	}
	return c.JSON(http.StatusOK, helpers.NewResponseSuccess("Task found", task))
}
