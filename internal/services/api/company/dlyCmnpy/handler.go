package dlyCmnpy

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Deiklov/diplom_backend/internal/common"
	"github.com/Deiklov/diplom_backend/internal/models"
	"github.com/Deiklov/diplom_backend/internal/services/api/company"
	diplom_backend "github.com/Deiklov/diplom_backend/internal/services/prediction/pb"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/Finnhub-Stock-API/finnhub-go"
	sdk "github.com/TinkoffCreditSystems/invest-openapi-go-sdk"
	"github.com/antihax/optional"
	"github.com/asaskevich/govalidator"
	"github.com/dgrijalva/jwt-go"
	"github.com/doug-martin/goqu/v9"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
	pbtime "google.golang.org/protobuf/types/known/timestamppb"
	"math/rand"
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
	tinkoffAPIURL string
	tinkoffToken  string
	tinkoffAPIWS  string
	sdkClient     *sdk.SandboxRestClient
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
		tinkoffAPIURL: "https://api-invest.tinkoff.ru/openapi/sandbox",
		tinkoffToken:  "t.UE-TeGMgnVeOVaoBYl7uk33-QtM9k2KwZTc7VyI1ubJErMxsVQvmYb92eRa157bm6XPjx74NGDIYfxSecNrdEQ",
		tinkoffAPIWS:  "wss://api-invest.tinkoff.ru/openapi/md/v1/md-openapi/ws",
		sdkClient:     sdk.NewSandboxRestClient("t.UE-TeGMgnVeOVaoBYl7uk33-QtM9k2KwZTc7VyI1ubJErMxsVQvmYb92eRa157bm6XPjx74NGDIYfxSecNrdEQ"),
	}
	mwareJWT := middleware.JWT([]byte("bc06c2d9-00cd-49e0-9f94-ef9257713803"))

	router.POST("/api/v1/company", handler.Create)
	router.GET("/api/v1/companies/favorite", handler.GetFavoriteList, mwareJWT)
	router.GET("/api/v1/companies", handler.GetAllCompaniesList)
	router.DELETE("/api/v1/company/favorite", handler.DeleteFavorite, mwareJWT)
	router.POST("/api/v1/company/favorite", handler.AddFavorite, mwareJWT)
	router.GET("/api/v1/company/page/:slug", handler.PersonalCompanyPage)
	router.GET("/api/v1/company/predict/:slug", handler.CompanyPredict)
	router.GET("/api/v1/companies/search/:slug", handler.CompanySearch)
	router.GET("/api/v1/market/candles", handler.GetCandles)
	router.GET("/api/v1/ws/market/candles/:slug", handler.GetRealTimeData)

}
func (usHttp *CompanyHttp) Create(ctx echo.Context) error {
	//curl --header "Content-Type: application/json" --request POST  --data '{"name":"ADBE"}'  http://localhost:8080/api/v1/company
	cmpny := models.Company{}
	if err := ctx.Bind(&cmpny); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	_, err := govalidator.ValidateStruct(cmpny)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: usHttp.FNapiKey,
	})
	profile2, resp, err := usHttp.finnhubClient.CompanyProfile2(auth, &finnhub.CompanyProfile2Opts{Symbol: optional.NewString(cmpny.Name)})
	//при плохом ответе длина будет 2, берем с запасом
	if resp == nil || resp.ContentLength != -1 {
		return echo.NewHTTPError(http.StatusUnprocessableEntity,
			"Stocks with such ticker doesn't exist")
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}
	cmpny = usHttp.FinhubProfileToModel(profile2)

	cmpnyFromDB, err := usHttp.UseCase.Create(cmpny)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, cmpnyFromDB)
}

func (usHttp *CompanyHttp) DeleteFavorite(ctx echo.Context) error {
	like := models.LikeUnlikeCompany{}
	if err := ctx.Bind(&like); err != nil {
		logger.Error(err)
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	userID := usHttp.GetUserID(ctx.Get("user").(*jwt.Token))
	err := usHttp.UseCase.DelFavorite(userID, models.Company{Ticker: like.Ticker})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.NoContent(http.StatusOK)

}

func (usHttp *CompanyHttp) AddFavorite(ctx echo.Context) error {
	like := models.LikeUnlikeCompany{}
	if err := ctx.Bind(&like); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	userID := usHttp.GetUserID(ctx.Get("user").(*jwt.Token))
	err := usHttp.UseCase.AddFavorite(userID, models.Company{
		Ticker: like.Ticker,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.NoContent(http.StatusOK)
}

//work
func (usHttp *CompanyHttp) GetFavoriteList(ctx echo.Context) error {
	cmpnys := []models.Company{}
	err := usHttp.dbsqlx.Select(&cmpnys, "SELECT cmp.id, ticker,name,ipo,country from company_by_users  cmpusers join companies cmp on cmp.id=cmpusers.company_id ")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, cmpnys)

}

//work
func (usHttp *CompanyHttp) GetAllCompaniesList(ctx echo.Context) error {
	companiesFromDB, err := usHttp.UseCase.GetAllCompanies()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, companiesFromDB)

}

//добавляется description и в будущем атрибуты при возврате инфы
func (usHttp *CompanyHttp) PersonalCompanyPage(ctx echo.Context) error {
	stocksSlug := ctx.Param("slug")
	if !govalidator.IsAlphanumeric(stocksSlug) {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid company slug!")
	}
	var cmpny models.Company
	err := usHttp.dbsqlx.Get(&cmpny, "SELECT id, name, description, country, attributes, ipo, ticker, logo, weburl from companies where ticker=$1", stocksSlug)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, cmpny)

}

func (usHttp *CompanyHttp) CompanyPredict(ctx echo.Context) error {
	ticker := ctx.Param("slug")
	query := diplom_backend.PredictionReq{
		StocksName: ticker,
		EndedTime:  pbtime.New(time.Now().Local().Add(2 * time.Hour)),
		Step:       2,
	}
	resp, err := usHttp.predictCL.Predict(ctx.Request().Context(), &query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return ctx.JSON(http.StatusOK, resp.TimeSeries)

}

//поиск только по slug
func (usHttp *CompanyHttp) CompanySearch(ctx echo.Context) error {
	stocksSlug := ctx.Param("slug")
	if !govalidator.IsAlphanumeric(stocksSlug) {
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

func (usHttp *CompanyHttp) GetCandles(ctx echo.Context) error {
	//todo добавить логику ограничений, по датам определяем минимальный интервал
	hub := sentryecho.GetHubFromContext(ctx)
	if hub == nil {
		logger.Info("nil sentry hub")
	}
	ticker := ctx.QueryParam("ticker")
	from := ctx.QueryParam("from")
	to := ctx.QueryParam("to")
	interval := ctx.QueryParam("interval")
	if interval == "" {
		interval = string(sdk.CandleInterval1Day)
	}
	fromTime, err := time.Parse(time.RFC3339, from)
	if err != nil {
		logger.Info(err)
		hub.CaptureMessage(err.Error())
		fromTime = time.Now().AddDate(-1, 0, 0)
	}
	toTime, err := time.Parse(time.RFC3339, to)
	if err != nil {
		logger.Info(err)
		hub.CaptureMessage(err.Error())
		toTime = time.Now()
	}

	instruments, err := usHttp.sdkClient.InstrumentByTicker(context.Background(), ticker)
	if err != nil || len(instruments) == 0 {
		if len(instruments) == 0 {
			err = errors.New("Zero length of instruments array")
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	candles, err := usHttp.sdkClient.Candles(context.Background(), fromTime, toTime, sdk.CandleInterval(interval), instruments[0].FIGI)
	if err != nil || len(candles) == 0 {
		if len(candles) == 0 {
			err = errors.New("Zero length of candles array")
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, candles)
}
func (usHttp *CompanyHttp) GetRealTimeData(ctx echo.Context) error {
	ticker := ctx.Param("slug")
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upgrader.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer func() { _ = ws.Close() }()

	instruments, err := usHttp.sdkClient.InstrumentByTicker(context.Background(), ticker)
	if err != nil || len(instruments) == 0 {
		if len(instruments) == 0 {
			err = errors.New("Zero length of instruments array")
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	client, err := sdk.NewStreamingClient(ctx.Logger(), usHttp.tinkoffToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	defer client.Close()

	go func() {
		if err := client.RunReadLoop(func(event interface{}) error {
			fmt.Println(event)
			if err := ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v", event.(sdk.OrderBookEvent)))); err != nil {
				return err
			}
			return nil
		}); err != nil {
			logger.Error(err)
		}
	}()
	if err := client.SubscribeOrderbook(instruments[0].FIGI, 1, requestID()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	time.Sleep(20 * time.Second)
	logger.Info("оптиска от стакана")
	if err := client.UnsubscribeOrderbook(instruments[0].FIGI, 1, requestID()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(http.StatusOK, "kek")
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// Генерируем уникальный ID для запроса
func requestID() string {
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
