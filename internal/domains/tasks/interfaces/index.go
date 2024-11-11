package interfaces

import (
	"github.com/supachai1998/task_services/internal/entities"
)

type TaskRepository interface {
	Create(task *entities.Task) error
	Update(task *entities.TaskUpdate) error
	GetByID(id uint) (*entities.Task, error)
	DeleteByID(id uint) error
	List() ([]entities.Task, error)
}
