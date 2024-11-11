package usecases

import (
	"github.com/supachai1998/task_services/internal/entities"
)

func (u *usecase) CreateTask(task *entities.Task) error {
	return u.taskRepo.Create(task)
}
