package entities

import "gorm.io/gorm"

type TaskStatus string

const (
	TaskStatusToDo       TaskStatus = "TO_DO"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusDone       TaskStatus = "DONE"
)

type Task struct {
	Id          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"not null;type:varchar(100)" json:"title"`
	Description string         `gorm:"not null;type:text" json:"description"`
	Status      TaskStatus     `gorm:"not null;default:TO_DO;" swagger:"enum(TO_DO,IN_PROGRESS,DONE)" json:"status"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Task) TableName() string {
	return TableNameTask
}

type TaskUpdate struct {
	Id          uint        `json:"id"`
	Title       *string     `json:"title"`
	Description *string     `json:"description"`
	Status      *TaskStatus `json:"status"`
}

func (TaskUpdate) TableName() string {
	return TableNameTask
}

// beforeCreate
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	t.Status = TaskStatusToDo
	return
}
