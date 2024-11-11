package repository

import (
	"github.com/supachai1998/task_services/internal/domains/tasks/interfaces"
	"github.com/supachai1998/task_services/internal/entities"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type repository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) interfaces.TaskRepository {
	return &repository{db}
}

func (r *repository) Create(task *entities.Task) error {
	return r.db.Create(task).Error
}

func (r *repository) Update(task *entities.TaskUpdate) error {
	return r.db.Clauses(clause.Returning{}).Where("id = ?", task.Id).Updates(task).Error
}

func (r *repository) GetByID(id uint) (*entities.Task, error) {
	var task entities.Task
	err := r.db.First(&task, id).Error
	return &task, err
}

func (r *repository) DeleteByID(id uint) error {
	return r.db.Delete(&entities.Task{}, id).Error
}

func (r *repository) List() ([]entities.Task, error) {
	var tasks []entities.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}
