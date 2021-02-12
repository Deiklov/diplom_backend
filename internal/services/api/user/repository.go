package user

import (
	"github.com/Deiklov/diplom_backend/internal/models"
)

//go:generate mockgen -source=repository.go -package=mocks -destination=./mocks/user_repo_mock.go
type Repository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Delete(id string) error
}
