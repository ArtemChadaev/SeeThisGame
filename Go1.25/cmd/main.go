package main

import (
	"os"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	repository2 "github.com/ArtemChadaev/SeeThisGame/internal/repository"
	"github.com/ArtemChadaev/SeeThisGame/internal/service"
	"github.com/ArtemChadaev/SeeThisGame/internal/transport/http"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("%s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("%s", err.Error())
	}
	db, err := repository2.NewPostgresDB(repository2.PostgresConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Database: viper.GetString("db.database"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("%s", err.Error())
	}
	redis, err := repository2.NewRedisClient(repository2.RedisConfig{
		Addr:     viper.GetString("redis.addr"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       viper.GetInt("redis.db"),
	})
	if err != nil {
		logrus.Fatalf("%s", err.Error())
	}
	repos := repository2.NewRepository(db)
	services := service.NewService(repos, redis)
	handlers := http.NewHandler(services, redis)

	srv := new(domain.rest)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error http: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
