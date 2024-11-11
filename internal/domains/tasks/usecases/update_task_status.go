package usecases

import (
	"errors"

	"github.com/samber/lo"
	"github.com/supachai1998/task_services/internal/entities"
)

// actionTransitions is a map of allowed transitions between task statuses.
var actionTransitions = map[entities.TaskStatus][]entities.TaskStatus{
	entities.TaskStatusToDo: {
		entities.TaskStatusInProgress,
		entities.TaskStatusDone,
	},
	entities.TaskStatusInProgress: {
		entities.TaskStatusDone,
	},
	entities.TaskStatusDone: {},
}

func (u *usecase) UpdateTaskStatus(id uint, status entities.TaskStatus) error {
	// Check if the status transition is allowed.
	if _, ok := actionTransitions[status]; !ok {
		return errors.New("invalid status")
	}
	return u.taskRepo.Update(&entities.TaskUpdate{
		Id:     id,
		Status: lo.ToPtr(status),
	})
}
