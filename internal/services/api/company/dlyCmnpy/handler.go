package dlyCmnpy

import (
	"context"
	"database/sql"
	"github.com/Deiklov/diplom_backend/internal/common"
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/company"
	diplom_backend "github.com/Deiklov/diplom_backend/internal/services/prediction/pb"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/doug-martin/goqu/v9"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	pbtime "google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"time"
)

type CompanyHttp struct {
	UseCase company.CompanyUCI
	DB      *sql.DB
	goquDb  *goqu.Database
	dbsqlx  *sqlx.DB
	common.UserGetter
	predictCL     diplom_backend.PredictAPIClient
	finnhubClient *finnhub.DefaultApiService
	FNapiKey      string
	common.CmpnyHelper
}

func AddRoutesWithHandler(router *echo.Echo, useCase company.CompanyUCI, db *sql.DB, client diplom_backend.PredictAPIClient, fnClient *finnhub.DefaultApiService, apiKey string) {
	handler := &CompanyHttp{
		UseCase:       useCase,
		DB:            db,
		dbsqlx:        sqlx.NewDb(db, "postgres"),
		goquDb:        goqu.New("postgres", db),
		predictCL:     client,
		finnhubClient: fnClient,
		FNapiKey:      apiKey,
	}
	mwareJWT := middleware.JWT([]byte("bc06c2d9-00cd-49e0-9f94-ef9257713803"))

	router.POST("/api/v1/company", handler.Create)
	router.GET("/api/v1/companies/favorite", handler.GetFavoriteList, mwareJWT)
	router.GET("/api/v1/companies", handler.GetAllCompaniesList)
	router.DELETE("/api/v1/company/favorite", handler.DeleteFavorite, mwareJWT)
	router.POST("/api/v1/company/favorite", handler.AddFavorite, mwareJWT)
	router.GET("/api/v1/company/page/:slug", handler.PersonalCompanyPage)
	router.GET("/api/v1/company/predict", handler.CompanyPredict, mwareJWT)
	router.GET("/api/v1/companies/search/:slug", handler.CompanySearch)

}
func (usHttp *CompanyHttp) Create(ctx echo.Context) error {
	//curl --header "Content-Type: application/json" --request POST  --data '{"name":"ADBE"}'  http://localhost:8080/api/v1/company
	cmpny := models.Company{}
	if err := ctx.Bind(&cmpny); err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	_, err := govalidator.ValidateStruct(cmpny)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: usHttp.FNapiKey,
	})
	profile2, resp, err := usHttp.finnhubClient.CompanyProfile2(auth, &finnhub.CompanyProfile2Opts{Symbol: optional.NewString(cmpny.Name)})
	//при плохом ответе длина будет 2, берем с запасом
	if resp == nil || resp.ContentLength != -1 {
		err = errors.New("Stocks with such ticker doesn't exist")
		logger.Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity, map[string]string{
			"error": err.Error(),
		})
	}
	cmpny = usHttp.FinhubProfileToModel(profile2)

	cmpnyFromDB, err := usHttp.UseCase.Create(cmpny)
	if err != nil {
		logger.Error(err, ctx.Request().Body)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, cmpnyFromDB)
}

func (usHttp *CompanyHttp) DeleteFavorite(ctx echo.Context) error {
	name := models.LikeUnlikeCompany{}
	if err := ctx.Bind(&name); err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	userID := usHttp.GetUserID(ctx.Get("user").(*jwt.Token))
	err := usHttp.UseCase.DelFavorite(userID, models.Company{Name: name.Name})
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.NoContent(http.StatusOK)

}

func (usHttp *CompanyHttp) AddFavorite(ctx echo.Context) error {
	name := models.LikeUnlikeCompany{}
	if err := ctx.Bind(&name); err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	userID := usHttp.GetUserID(ctx.Get("user").(*jwt.Token))
	err := usHttp.UseCase.AddFavorite(userID, models.Company{
		Name: name.Name,
	})
	if err != nil {
		logger.Error(err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.NoContent(http.StatusOK)
}

//work
func (usHttp *CompanyHttp) GetFavoriteList(ctx echo.Context) error {
	cmpnys := []models.Company{}
	err := usHttp.dbsqlx.Select(&cmpnys, "SELECT cmp.id, name,year,country from company_by_users  cmpusers join companies cmp on cmp.id=cmpusers.company_id ")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, cmpnys)

}

//work
func (usHttp *CompanyHttp) GetAllCompaniesList(ctx echo.Context) error {
	companiesFromDB, err := usHttp.UseCase.GetAllCompanies()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(200, companiesFromDB)

}

//добавляется description и в будущем атрибуты при возврате инфы
func (usHttp *CompanyHttp) PersonalCompanyPage(ctx echo.Context) error {
	stocksSlug := ctx.Param("slug")
	if !govalidator.IsBase64(stocksSlug) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid company slug",
		})
	}
	var cmpny models.Company
	err := usHttp.dbsqlx.Get(&cmpny, "SELECT id, name, description, country, attributes, ipo, ticker, logo, weburl from companies where ticker=$1", stocksSlug)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, cmpny)

}

func (usHttp *CompanyHttp) CompanyPredict(ctx echo.Context) error {
	query := diplom_backend.PredictionReq{
		StocksName: "AAPL",
		EndedTime:  pbtime.New(time.Now().Local().Add(2 * time.Hour)),
		Step:       2,
	}
	resp, err := usHttp.predictCL.Predict(ctx.Request().Context(), &query)
	if err != nil {
		logger.Error(err)
		return err
	}
	return ctx.JSON(http.StatusOK, resp.TimeSeries)

}

//поиск только по slug
func (usHttp *CompanyHttp) CompanySearch(ctx echo.Context) error {
	stocksSlug := ctx.Param("slug")
	if !govalidator.IsBase64(stocksSlug) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid company slug",
		})
	}
	companyFromDB, err := usHttp.UseCase.GetCompany(stocksSlug)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return ctx.JSON(http.StatusOK, companyFromDB)

}
