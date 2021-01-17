package http

import (
	"github.com/Deiklov/diplom_backend/internal/services/api/user"
	"github.com/labstack/echo"
)

type UserHandler struct {
	UseCase user.UseCase
}

func CreateHandler(router *echo.Echo, useCase user.UseCase, mw *middleware.Middleware) {
	handler := &UserHandler{
		UseCase: useCase,
	}
	router.POST("/api/settings", handler.Create, mw.Sanitize)
	router.GET("/api/profile/:id_or_nickname", handler.Get)
	router.GET("/api/settings", handler.GetAll, mw.CheckAuth, mw.CSRFmiddle) // получ все настройки
	router.PUT("/api/settings", handler.Update, mw.CheckAuth, mw.CSRFmiddle)
	router.DELETE("/api/settings", handler.Delete, mw.CheckAuth, mw.CSRFmiddle)
	//router.GET("/search/profile", handler.GetUsersByNicknamePart, mw.CheckAuth)
}
