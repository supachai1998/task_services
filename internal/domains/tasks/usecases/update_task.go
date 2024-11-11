package usecases

import (
	"errors"

	"github.com/supachai1998/task_services/internal/entities"
)

// actionTaskStatus is a map not allowed update task information.
var actionNotAllowed = map[entities.TaskStatus]bool{
	entities.TaskStatusDone: true,
}

func (u *usecase) UpdateTask(task *entities.TaskUpdate) error {
	currentTask, err := u.taskRepo.GetByID(task.Id)
	if err != nil {
		return err
	}
	if currentTask != nil {
		if _, ok := actionNotAllowed[currentTask.Status]; ok {
			return errors.New("this task status is done, cannot update")
		}
	}

	return u.taskRepo.Update(task)
}
