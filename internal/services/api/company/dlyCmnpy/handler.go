package dlyCmnpy

import (
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/company"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/asaskevich/govalidator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type CompanyHttp struct {
	UseCase company.CompanyUCI
}

func AddRoutesWithHandler(router *echo.Echo, useCase company.CompanyUCI) {
	handler := &CompanyHttp{
		UseCase: useCase,
	}
	router.POST("/api/v1/company", handler.Create)
	router.GET("/api/v1/user/:id_or_nickname", handler.Get, middleware.JWT([]byte("secret")))
	router.PUT("/api/v1/user", handler.Update, middleware.JWT([]byte("secret")))
	router.DELETE("/api/v1/user", handler.Delete, middleware.JWT([]byte("secret")))
	router.GET("/api/v1/whoami", handler.Get, middleware.JWT([]byte("secret")))
	router.POST("/api/v1/login", handler.Login, middleware.JWT([]byte("secret")))

}
func (usHttp *CompanyHttp) Create(ctx echo.Context) error {
	cmpny := models.Company{}
	if err := ctx.Bind(&cmpny); err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, "Can't parse company data")
	}
	_, err := govalidator.ValidateStruct(cmpny)
	if err != nil {
		logger.Error(err)
		return ctx.String(http.StatusBadRequest, err.Error())
	}
	usHttp.UseCase.Create()
	return nil
}
