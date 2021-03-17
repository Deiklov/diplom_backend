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
	mwareJWT := middleware.JWT([]byte("bc06c2d9-00cd-49e0-9f94-ef9257713803"))

	router.POST("/api/v1/company", handler.Create)
	router.GET("/api/v1/companies/favorite", handler.GetFavoriteList, mwareJWT)
	router.GET("/api/v1/companies", handler.GetAllCompaniesList)
	router.DELETE("/api/v1/company/favorite", handler.DeleteFavorite, mwareJWT)
	router.POST("/api/v1/company/favorite", handler.AddFavorite, mwareJWT)
	router.GET("/api/v1/company/favorite", handler.PersonalCompanyPage)
	router.GET("/api/v1/company/predict", handler.CompanyPredict, mwareJWT)

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
	cmpnyFromDB, err := usHttp.UseCase.Create(cmpny)
	if err != nil {
		logger.Error(err, ctx.Request().Body)
	}
	return ctx.JSON(http.StatusOK, cmpnyFromDB)
}

func (usHttp *CompanyHttp) DeleteFavorite(ctx echo.Context) error {
	return ctx.String(200, "del favorite ")

}

func (usHttp *CompanyHttp) AddFavorite(ctx echo.Context) error {
	return ctx.String(200, "add favorite")

}

func (usHttp *CompanyHttp) GetFavoriteList(ctx echo.Context) error {
	return ctx.String(200, "favorite list")

}

func (usHttp *CompanyHttp) GetAllCompaniesList(ctx echo.Context) error {
	return ctx.String(200, "all companies list")
}

func (usHttp *CompanyHttp) PersonalCompanyPage(ctx echo.Context) error {
	return ctx.String(200, "full info company")

}

func (usHttp *CompanyHttp) CompanyPredict(ctx echo.Context) error {
	return ctx.String(200, "predict company")

}
