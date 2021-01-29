package user

import (
	"github.com/Deiklov/diplom_backend/internal/models"
)

type UseCase interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	Delete(id string) error
}

