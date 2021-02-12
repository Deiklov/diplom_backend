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
	router.POST("/api/v1/user", handler.Create)
	router.GET("/api/v1/user/:id_or_nickname", handler.Get, middleware.JWT([]byte("secret")))
	router.PUT("/api/v1/user", handler.Update, middleware.JWT([]byte("secret")))
	router.DELETE("/api/v1/user", handler.Delete, middleware.JWT([]byte("secret")))
	router.GET("/api/v1/whoami", handler.Get, middleware.JWT([]byte("secret")))
	router.POST("/api/v1/login", handler.Login, middleware.JWT([]byte("secret")))

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
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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

func (usHttp *UserHttp) Get(ctx echo.Context) error {

	userCurr := ctx.Get("user").(*jwt.Token)
	claims := userCurr.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return ctx.String(http.StatusOK, "kek "+name)
}

func (usHttp *UserHttp) GetAll(ctx echo.Context) error {

	return ctx.String(http.StatusOK, "kek")
}

func (usHttp *UserHttp) Update(ctx echo.Context) error {

	return ctx.NoContent(http.StatusOK)
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
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}
	return t, nil
}
