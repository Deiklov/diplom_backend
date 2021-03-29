package common

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserGetter struct {
}

func (user *UserGetter) GetUserID(token *jwt.Token) string {
	claims := token.Claims.(jwt.MapClaims)
	userClaims := claims["user"].(map[string]interface{})
	return userClaims["id"].(string)
}

func (user *UserGetter) GetUser(token *jwt.Token) (*models.User, error) {
	claims := token.Claims.(jwt.MapClaims)
	userClaims := claims["user"].(map[string]interface{})
	createdAt, err := time.Parse(time.RFC3339, userClaims["created_at"].(string))
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, userClaims["updated_at"].(string))
	if err != nil {
		return nil, err
	}
	userData := &models.User{
		ID:        userClaims["id"].(string),
		Name:      userClaims["name"].(string),
		Email:     userClaims["email"].(string),
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
		DeletedAt: nil,
		Password:  "",
	}
	return userData, nil
}
