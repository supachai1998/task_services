package models

type UpdateTaskStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=IN_PROGRESS DONE"`
}

type CreateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100" example:"Later is never"`
	Description string `json:"description" validate:"required,min=3,max=25500" example:"When 'later' turns into 'never', it's just your code's way of saying it loves the TODO comments."`
}
type UpdateTaskRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=100" example:"Code runs, coffee fuels"`
	Description string `json:"description" validate:"required,min=3,max=25500" example:"Coding without coffee is like debugging without a console log."`
}
