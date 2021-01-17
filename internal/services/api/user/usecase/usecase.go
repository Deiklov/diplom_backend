package usecase

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"mime/multipart"
)

type UserUseCase struct {
	//sessionRepo session.Repository
	userRepo user.Repository
}

func CreateUseCase(userRepo_ user.Repository) user.UseCase {
	return &UserUseCase{
		userRepo: userRepo_,
	}
}

func (userUseCase *UserUseCase) Create(user *models.User, sessionExpires int32) (string, error) {
	err := userUseCase.userRepo.Create(user)
	if err != nil {
		logger.Error(err)
		return "", err
	}

	return "", nil
}

func (userUseCase *UserUseCase) GetByID(uid uint) (*models.User, error) {
	usr, err := userUseCase.userRepo.GetByID(uid)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usr, nil
}

func (userUseCase *UserUseCase) GetByNickname(nickname string) (*models.User, error) {
	usr, err := userUseCase.userRepo.GetByNickname(nickname)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return usr, nil
}

func (userUseCase *UserUseCase) GetUsersByNicknamePart(nicknamePart string, limit uint) ([]models.User, error) {
	users, err := userUseCase.userRepo.GetUsersByNicknamePart(nicknamePart, limit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return users, nil
}



func (userUseCase *UserUseCase) Delete(uid uint, sid string) error {
	return userUseCase.userRepo.Delete(uid)
}
