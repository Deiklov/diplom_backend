package repository

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
	errOwn "github.com/Deiklov/diplom_backend/pkg/errors"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type UserStore struct {
	DB *gorm.DB
}

func CreateRepository(db *gorm.DB) user.Repository {
	return &UserStore{DB: db}
}

func (userStore *UserStore) Create(usr *models.User) error {
	if err := userStore.DB.Create(usr).Error; err != nil {
		logger.Error(err)
		return errOwn.ErrConflict
	}
	return nil
}

func (userStore *UserStore) GetByID(id uint) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.DB.Where("id = ?", id).First(&usr).Error; err != nil {
		logger.Error(err)
		return nil, errOwn.ErrUserNotFound
	}
	return usr, nil
}

func (userStore *UserStore) GetByNickname(nickname string) (*models.User, error) {
	usr := new(models.User)
	if err := userStore.DB.Where("nickname = ?", nickname).First(&usr).Error; err != nil {
		logger.Error(err)
		return nil, errOwn.ErrUserNotFound
	}
	return usr, nil
}

func (userStore *UserStore) Delete(id uint) error {
	if err := userStore.DB.Where("id = ?", id).Delete(models.User{}).Error; err != nil {
		logger.Error(err)
		return errOwn.ErrUserNotFound
	}
	return nil
}

func (userStore *UserStore) GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error) {
	var users []models.User
	err := userStore.DB.Limit(limit).Where("nickname LIKE ?", nicknamePart+"%").Find(&users).Error
	if err != nil {
		logger.Error(err)
		return nil, errOwn.ErrUserNotFound
	}
	return users, nil
}
