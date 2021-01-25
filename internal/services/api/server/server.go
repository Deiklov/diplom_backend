package server

import (
	"github.com/Deiklov/diplom_backend/config"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Server struct {
	IP   string
	Port uint
	Conf config.Config
}

func NewServer(ip string, port uint) *Server {
	//config.ReadConfig()
	return &Server{
		IP:   ip,
		Port: port,
		Conf: config.C,
	}
}

func (serv *Server) Run() {
	//config.ReadConfig()

	//connectionString := fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d",
	//	serv.Conf.Database.DBName, serv.Conf.Database.User, serv.Conf.Database.Password, serv.IP, 5432)
	//_, err := sqlx.Connect("pgx", connectionString)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//_, err = pgxpool.Connect(context.Background(), connectionString)
	//
	//if err != nil {
	//	panic(err)
	//}
	router := echo.New()
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
