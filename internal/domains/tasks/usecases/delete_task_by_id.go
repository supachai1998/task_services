package usecases

func (u *usecase) DeleteTaskByID(id uint) error {
	return u.taskRepo.DeleteByID(id)
}
