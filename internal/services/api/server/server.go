package server

import (
	goSQL "database/sql"
	"fmt"
	"github.com/Deiklov/diplom_backend/config"
	httpUser "github.com/Deiklov/diplom_backend/internal/services/api/user/delivery/http"
	"github.com/Deiklov/diplom_backend/internal/services/api/user/repUser"
	"github.com/Deiklov/diplom_backend/internal/services/api/user/ucUser"
	"net/http"

	//"github.com/doug-martin/goqu/v9"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
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
			"http://bmstu-romanov.xyz", "https://localhost:3000", "http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut,
			http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     nil,
		AllowCredentials: true,
		ExposeHeaders:    nil,
		MaxAge:           0,
	}))
	//router.Use(middleware.CSRFWithConfig(middleware.DefaultCSRFConfig))
	router.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	router.Use(middleware.Recover())

	userRepo := repUser.CreateRepository(pdb)
	userUC := ucUser.CreateUseCase(userRepo)
	httpUser.AddRoutesWithHandler(router, userUC)

	router.Logger.Fatal(router.Start(":8080"))

}
