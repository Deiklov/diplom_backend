package http

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"time"
)

type UserHttp struct {
	UseCase user.UseCase
}

func AddRoutesWithHandler(router *echo.Echo, useCase user.UseCase) {
	handler := &UserHttp{
		UseCase: useCase,
	}
	mwareJWT := middleware.JWT([]byte("bc06c2d9-00cd-49e0-9f94-ef9257713803"))
	router.POST("/api/v1/user", handler.Create)
	router.GET("/api/v1/user/:id_or_nickname", handler.Get, mwareJWT)
	router.PUT("/api/v1/user", handler.Update, mwareJWT)
	router.DELETE("/api/v1/user", handler.Delete, mwareJWT)
	router.GET("/api/v1/user", handler.Get, mwareJWT)
	router.POST("/api/v1/login", handler.Login)

}
func (usHttp *UserHttp) Create(ctx echo.Context) error {
	usr := models.User{}
	if err := ctx.Bind(&usr); err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Can't parse user data")
	}
	_, err := govalidator.ValidateStruct(usr)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	if err := usHttp.UseCase.Create(&models.User{
		ID:       uuid.NewString(),
		Name:     usr.Name,
		Email:    usr.Email,
		Password: usr.Password,
	}); err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	token, err := usHttp.login(&models.AuthData{
		Email:    usr.Email,
		Password: usr.Password,
	})
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
func (usHttp *UserHttp) parseUser(userClaims map[string]interface{}) (*models.User, error) {
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
func (usHttp *UserHttp) Get(ctx echo.Context) error {

	userCurr := ctx.Get("user").(*jwt.Token)
	claims := userCurr.Claims.(jwt.MapClaims)
	userClaims := claims["user"].(map[string]interface{})
	userData, err := usHttp.parseUser(userClaims)
	if err != nil {
		return ctx.String(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, userData)
}

func (usHttp *UserHttp) GetAll(ctx echo.Context) error {

	return ctx.String(http.StatusOK, "kek")
}

func (usHttp *UserHttp) Update(ctx echo.Context) error {
	usr := models.User{}
	if err := ctx.Bind(&usr); err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Can't parse user data")
	}
	_, err := govalidator.ValidateStruct(usr)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	userFromDB, err := usHttp.UseCase.Update(&usr)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, userFromDB)
}

func (usHttp *UserHttp) Delete(ctx echo.Context) error {

	return ctx.NoContent(http.StatusOK)
}

func (usHttp *UserHttp) Login(ctx echo.Context) error {
	authData := &models.AuthData{}
	if err := ctx.Bind(authData); err != nil {
		return err
	}
	_, err := govalidator.ValidateStruct(authData)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	token, err := usHttp.login(authData)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

func (usHttp *UserHttp) login(data *models.AuthData) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	usr, err := usHttp.UseCase.GetByEmail(data.Email)
	if err != nil {
		return "", err
	}
	claims["user"] = usr
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte("bc06c2d9-00cd-49e0-9f94-ef9257713803"))
	if err != nil {
		return "", err
	}
	return t, nil
}
