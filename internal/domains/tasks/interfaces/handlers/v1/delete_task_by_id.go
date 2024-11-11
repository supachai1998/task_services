package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/supachai1998/task_services/internal/helpers"
)

// UpdateTask updates task details
// @Summary Update task details
// @Description Update a task by its unique ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Success 204 "Task deleted successfully"
// @Failure 400 {object} models.ResponseError "Invalid input"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/tasks/{id} [delete]
func (h *Handler) DeleteTaskByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, helpers.NewResponseError("Invalid ID format", "error"))
	}
	if err := h.TaskUsecase.DeleteTaskByID(uint(id)); err != nil {
		return c.JSON(http.StatusNotFound, helpers.NewResponseError("Task not found", "error"))
	}
	return c.NoContent(http.StatusNoContent)
}
