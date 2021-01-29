package http

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/middleware"
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type UserHttp struct {
	UseCase user.UseCase
}

func AddRoutesWithHandler(router *echo.Echo, useCase user.UseCase, mw *middleware.Middleware) {
	handler := &UserHttp{
		UseCase: useCase,
	}
	router.POST("/api/user", handler.Create)
	router.GET("/api/user/:id_or_nickname", handler.Get)
	router.PUT("/api/user", handler.Update)
	router.DELETE("/api/user", handler.Delete)

}
func (usHttp *UserHttp) Create(ctx echo.Context) error {
	usr := models.User{}
	if err := ctx.Bind(usr); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Can't parse user data")
	}
	if err := usHttp.UseCase.Create(&models.User{
		ID:    uuid.Must(uuid.NewV4()).String(),
		Name:  usr.Name,
		Phone: usr.Phone,
	}); err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return ctx.NoContent(http.StatusOK)
}

func (usHttp *UserHttp) Get(ctx echo.Context) error {

	return ctx.String(http.StatusOK, "kek")
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
