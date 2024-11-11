package usecases

import "github.com/supachai1998/task_services/internal/entities"

func (u *usecase) GetTaskByID(id uint) (*entities.Task, error) {
	return u.taskRepo.GetByID(id)
}
