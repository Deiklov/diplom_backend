package server

import (
	"fmt"
	"github.com/Deiklov/diplom_backend/config"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
	"golang.org/x/net/context"
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
	config.ReadConfig()

	conn, err := pgx.ConnectConfig(context.Background(), &pgx.ConnConfig{
		Config: pgconn.Config{
			Host:            serv.IP,
			Port:            5432,
			Database:        serv.Conf.Database.DBName,
			User:            serv.Conf.Database.User,
			Password:        serv.Conf.Database.Password,
			TLSConfig:       nil,
			DialFunc:        nil,
			LookupFunc:      nil,
			BuildFrontend:   nil,
			RuntimeParams:   nil,
			Fallbacks:       nil,
			ValidateConnect: nil,
			AfterConnect:    nil,
			OnNotice:        nil,
			OnNotification:  nil,
		},
		Logger:               nil,
		LogLevel:             0,
		BuildStatementCache:  nil,
		PreferSimpleProtocol: false,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
}
