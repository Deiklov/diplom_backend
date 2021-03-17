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
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"net/http"

	//"github.com/doug-martin/goqu/v9"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4/middleware"
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
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		serv.Conf.Database.User, serv.Conf.Database.Password, serv.Conf.Database.DBName)
	pdb, err := goSQL.Open("pgx", connectionString)
	if err != nil {
		panic(err)
	}

	router := echo.New()
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper: middleware.DefaultSkipper,
		AllowOrigins: []string{"https://localhost:80", "http://localhost:80", "https://bmstu-romanov.xyz",
			"http://bmstu-romanov.xyz", "https://localhost:3000", "http://localhost:3000",
			"http://23.111.206.228:80", "https://23.111.206.228:80", "http://23.111.206.228:8080", "https://23.111.206.228:8080"},
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

	//if err := sentry.Init(sentry.ClientOptions{
	//	Dsn: "https://325b0bf7ab154d349f3e51010d1964e7@o553507.ingest.sentry.io/5680876",
	//}); err != nil {
	//	logger.Fatal(err)
	//}
	//router.Use(sentryecho.New(sentryecho.Options{}))
	//sentry.CaptureMessage("It works!")
	//router.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	router.Use(middleware.Recover())
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(router)

	userRepo := repUser.CreateRepository(pdb)
	userUC := ucUser.CreateUseCase(userRepo)
	httpUser.AddRoutesWithHandler(router, userUC) //добавит юзерские роуты

	cmpnyRepo := repCmpny.CreateRepCmpny(pdb)
	cmpnyUCase := ucCmnpy.CreateUseCase(cmpnyRepo)
	dlyCmnpy.AddRoutesWithHandler(router, cmpnyUCase) //добавит роуты компании

	router.Logger.Fatal(router.Start(":8080"))

}
