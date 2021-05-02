package server

import (
	goSQL "database/sql"
	"fmt"
	"github.com/Deiklov/diplom_backend/config"
	"github.com/Deiklov/diplom_backend/internal/services/api/company/dlyCmnpy"
	"github.com/Deiklov/diplom_backend/internal/services/api/company/repCmpny"
	"github.com/Deiklov/diplom_backend/internal/services/api/company/ucCmnpy"
	httpUser "github.com/Deiklov/diplom_backend/internal/services/api/user/delivery/http"
	"github.com/Deiklov/diplom_backend/internal/services/api/user/repUser"
	"github.com/Deiklov/diplom_backend/internal/services/api/user/ucUser"
	diplom_backend "github.com/Deiklov/diplom_backend/internal/services/prediction/pb"
	"github.com/Deiklov/diplom_backend/pkg/logger"
	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
)

type Server struct {
	IP   string
	Port uint
	Conf config.Config
}

func NewServer(ip string, port uint) *Server {
	config.ReadConfig()
	return &Server{
		IP:   ip,
		Port: port,
		Conf: config.C,
	}
}

func (serv *Server) Run() {
	fmt.Println(func() string { dir, _ := os.Getwd(); return dir }())
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s",
		serv.Conf.Database.User, serv.Conf.Database.Password, serv.Conf.Database.Addr, serv.Conf.Database.DBName)
	pdb, err := goSQL.Open("pgx", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	router := echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"https://bmstu-romanov.xyz", "http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut,
			http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     nil,
		AllowCredentials: true,
		ExposeHeaders:    nil,
		MaxAge:           0,
	}))
	//router.Use(middleware.CSRFWithConfig(middleware.DefaultCSRFConfig))
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_custom}","id":"${id}","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency_human":"${latency_human}"` +
			`,"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	router.Use(middleware.Recover())

	if err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://325b0bf7ab154d349f3e51010d1964e7@o553507.ingest.sentry.io/5680876",
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
	router.Use(sentryecho.New(sentryecho.Options{}))

	router.HTTPErrorHandler = func(err error, c echo.Context) {
		he, ok := err.(*echo.HTTPError)
		if ok {
			if he.Internal != nil {
				if herr, ok := he.Internal.(*echo.HTTPError); ok {
					he = herr
				}
			}
		} else {
			he = &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: http.StatusText(http.StatusInternalServerError),
			}
		}

		// Issue #1426
		code := he.Code
		message := he.Message
		if router.Debug {
			message = err.Error()
		} else if m, ok := message.(string); ok {
			message = echo.Map{"error": m}
		}

		// Send response
		if !c.Response().Committed {
			//логируем в консоль и в сентри
			logger.Error(err)
			if hub := sentryecho.GetHubFromContext(c); hub != nil {
				//используем message тк call stack будет в ошибке и меньше длина логов
				hub.CaptureMessage(err.Error())
			}
			if c.Request().Method == http.MethodHead { // Issue #608
				err = c.NoContent(he.Code)
			} else {
				err = c.JSON(code, message)
			}
			if err != nil {
				router.Logger.Error(err)
			}
		}
	}

	p := prometheus.NewPrometheus("echo", nil)
	p.Use(router)

	addr := "localhost:9999"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		logger.Warn(err)
	}
	defer conn.Close()
	client := diplom_backend.NewPredictAPIClient(conn)

	finnhubClient := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi

	userRepo := repUser.CreateRepository(pdb)
	userUC := ucUser.CreateUseCase(userRepo)
	httpUser.AddRoutesWithHandler(router, userUC) //добавит юзерские роуты

	cmpnyRepo := repCmpny.CreateRepCmpny(pdb)
	cmpnyUCase := ucCmnpy.CreateUseCase(cmpnyRepo)
	//принимает репозиторий, чтобы быстрее шла разработка
	dlyCmnpy.AddRoutesWithHandler(router, cmpnyUCase, pdb, client, finnhubClient, "c0ilbh748v6ot9ddgc0g")

	router.Logger.Fatal(router.Start(":8080"))

}
