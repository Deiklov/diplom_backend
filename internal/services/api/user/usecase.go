package user

import (
	"github.com/Deiklov/diplom_backend/internal/models"
)

type UseCase interface {
	Create(user *models.User, sessionExpires int32) (string, error)
	GetByID(uid uint) (*models.User, error)
	GetByNickname(nickname string) (*models.User, error)
	GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error)
	Delete(uid uint, sid string) error
}