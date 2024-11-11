package usecases

import (
	"github.com/supachai1998/task_services/internal/domains/tasks/interfaces"
	"github.com/supachai1998/task_services/internal/entities"
)

type TaskUsecase interface {
	CreateTask(task *entities.Task) error
	UpdateTask(task *entities.TaskUpdate) error
	UpdateTaskStatus(id uint, status entities.TaskStatus) error
	GetTaskByID(id uint) (*entities.Task, error)
	DeleteTaskByID(id uint) error
	ListTasks() ([]entities.Task, error)
}

type usecase struct {
	taskRepo interfaces.TaskRepository
}

func NewTaskUsecase(taskRepo interfaces.TaskRepository) TaskUsecase {
	return &usecase{taskRepo}
}
