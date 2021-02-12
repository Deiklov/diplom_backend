package ucUser

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
	"golang.org/x/crypto/bcrypt"
)

func CreateUseCase(userRepo_ user.Repository) user.UseCase {
	return &userUCase{
		userRepo: userRepo_,
	}
}

type userUCase struct {
	userRepo user.Repository
}

func (uc *userUCase) Create(user *models.User) error {
	pass, _ := uc.HashPassword(user.Password)
	user.Password = pass
	return uc.userRepo.Create(user)
}
func (uc *userUCase) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func (uc *userUCase) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func (uc *userUCase) GetByEmail(email string) (*models.User, error) {
	return uc.userRepo.GetByEmail(email)
}

func (uc *userUCase) GetByID(uid string) (*models.User, error) {
	return uc.userRepo.GetByID(uid)
}

func (uc *userUCase) Delete(uid string) error {
	return uc.userRepo.Delete(uid)
}
