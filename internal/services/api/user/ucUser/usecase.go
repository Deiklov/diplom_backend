package ucUser

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
)



func CreateUseCase(userRepo_ user.Repository) user.UseCase {
	return &userUCase{
		userRepo: userRepo_,
	}
}
type userUCase struct{
	userRepo user.Repository
}

func (userUCase) Create(user *models.User) error {
	panic("implement me")
}

func (userUCase) GetByID(uid string) (*models.User, error) {
	panic("implement me")
}

func (userUCase) Delete(uid string) error {
	panic("implement me")
}
