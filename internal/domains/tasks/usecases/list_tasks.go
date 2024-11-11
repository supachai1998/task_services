package usecases

import "github.com/supachai1998/task_services/internal/entities"

func (u *usecase) ListTasks() ([]entities.Task, error) {
	return u.taskRepo.List()
}
