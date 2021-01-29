package server

import (
	goSQL "database/sql"
	"fmt"
	"github.com/Deiklov/diplom_backend/config"
	mwareInternal "github.com/Deiklov/diplom_backend/internal/services/api/middleware"
	httpUser "github.com/Deiklov/diplom_backend/internal/services/api/user/delivery/http"
	"github.com/Deiklov/diplom_backend/internal/services/api/user/repUser"
	"github.com/Deiklov/diplom_backend/internal/services/api/user/ucUser"
	"github.com/bxcodec/faker/v3"
	"github.com/doug-martin/goqu/v9"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/echo/v4"
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
	var goquDb *goqu.Database
	connectionString := fmt.Sprintf("postgres://%s:%s@localhost:5432/%s?sslmode=disable",
		serv.Conf.Database.User, serv.Conf.Database.Password, serv.Conf.Database.DBName)
	pdb, err := goSQL.Open("pgx", connectionString)
	if err != nil {
		panic(err)
	}
	goquDb = goqu.New("postgres", pdb)
	sqlResult, err := goquDb.Insert("users").Cols("id", "phone").
		Vals(goqu.Vals{faker.UUIDHyphenated(), faker.Phonenumber()}).Executor().Exec()

	fmt.Println(sqlResult)
	if err != nil {
		fmt.Println(err)
		return
	}

	router := echo.New()
	router.Use(middleware.JWTWithConfig(middleware.DefaultJWTConfig))
	router.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	router.Use(middleware.CSRFWithConfig(middleware.DefaultCSRFConfig))
	router.Use(middleware.LoggerWithConfig(middleware.DefaultLoggerConfig))
	router.Use(middleware.BodyLimit("6MB"))
	router.Use(middleware.Recover())
	router.Use(middleware.Secure())
	userRepo := repUser.CreateRepository(pdb)
	userUC := ucUser.CreateUseCase(userRepo)
	mware := mwareInternal.CreateMiddleware(userUC)
	httpUser.AddRoutesWithHandler(router, userUC, mware)

	router.GET("/api/kek", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! kek")
	})
	router.GET("/api/kek/:shpek", func(c echo.Context) error {
		return c.String(http.StatusCreated, "Hello, World! shpek")
	})
	router.GET("/api/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World! api ")
	})
	router.GET("/api/mda", func(c echo.Context) error {
		return c.String(http.StatusNoContent, "mdamdamda ")
	})
	router.Logger.Fatal(router.Start(":8080"))

}
