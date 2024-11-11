package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/supachai1998/task_services/internal/helpers"
)

// ListTasks handles task listing
// @Summary List all tasks
// @Description List all tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseSuccess{data=[]entities.Task} "Tasks listed successfully"
// @Failure 500 {object} models.ResponseError "Internal server error"
// @Router /v1/tasks [get]
func (h *Handler) ListTasks(c echo.Context) error {
	tasks, err := h.TaskUsecase.ListTasks()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, helpers.NewResponseError(err.Error(), "error"))
	}
	return c.JSON(http.StatusOK, helpers.NewResponseSuccess("Tasks listed", tasks))
}
