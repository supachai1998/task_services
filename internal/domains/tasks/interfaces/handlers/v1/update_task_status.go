package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/supachai1998/task_services/internal/domains/tasks/models"
	"github.com/supachai1998/task_services/internal/entities"
	"github.com/supachai1998/task_services/internal/helpers"
)

// UpdateTask updates task details
// @Summary Update task details
// @Description Update a task by its unique ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param body body models.UpdateTaskStatusRequest true "Task details"
// @Success 200 {object} models.ResponseSuccess{} "Task updated successfully"
// @Failure 400 {object} models.ResponseError "Invalid input"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/tasks/{id}/status [patch]
func (h *Handler) UpdateTaskStatus(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.NewResponseError("Invalid ID format", "error"))
	}
	var req models.UpdateTaskStatusRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.NewResponseError(err.Error(), "error"))
	}

	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, helpers.NewResponseError(err.Error(), "error"))
	}

	if err := h.TaskUsecase.UpdateTaskStatus(uint(id), entities.TaskStatus(req.Status)); err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.NewResponseError(err.Error(), "error"))
	}
	return c.JSON(http.StatusOK, helpers.NewResponseSuccess("Task status updated", ""))
}
