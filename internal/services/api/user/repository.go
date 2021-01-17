package user

import (
	"github.com/Deiklov/diplom_backend/internal/models"
)

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/user_repo_mock.go
type Repository interface {
	Create(user *models.User) error
	GetByID(id uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	Delete(id uint) error
	GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error)
}
