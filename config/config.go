package config

import (
	"fmt"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"
)

type Config struct {
	Database struct {
		User                 string
		Password             string
		Net                  string
		Addr                 string
		DBName               string
		AllowNativePasswords bool
		Params               struct {
			ParseTime string
		}
	}
	Server struct {
		Address string
	}
}

var C Config

func ReadConfig() {
	Config := &C

	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("config/")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}

	if err := viper.Unmarshal(&Config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//err := godotenv.Load("main.env")
	//if err != nil {
	//	log.Fatal(err)
	//}
	Config.Database.DBName = os.Getenv("POSTGRES_DB")
	Config.Database.Password = os.Getenv("POSTGRES_PASSWORD")
	Config.Database.User = os.Getenv("POSTGRES_USER")
	Config.Database.Addr = os.Getenv("POSTGRES_HOST")

	spew.Dump(C)
}
