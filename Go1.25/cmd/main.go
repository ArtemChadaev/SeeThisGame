package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ArtemChadaev/SeeThisGame/internal/domain"
	"github.com/ArtemChadaev/SeeThisGame/internal/repository"
	"github.com/ArtemChadaev/SeeThisGame/internal/service"
	"github.com/ArtemChadaev/SeeThisGame/internal/transport/rest"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// 1. Настройка логгера
	logrus.SetFormatter(new(logrus.JSONFormatter))

	// 2. Инициализация конфигурации
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Warn("no .env file found, using environment variables")
	}

	// 3. Подключение к БД (Postgres)
	db, err := repository.NewPostgresDB(repository.PostgresConfig{
		// Теперь viper будет проверять переменные окружения, если мы настроим его ниже
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Database: viper.GetString("db.database"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	}

	// ЗАПУСК МИГРАЦИЙ
	logrus.Info("Running database migrations...")
	if err := repository.RunMigrations(db); err != nil {
		logrus.Fatalf("Migrations failed: %s", err.Error())
	}
	logrus.Info("Migrations applied successfully!")

	// 4. Подключение к Redis
	redisClient, err := repository.NewRedisClient(repository.RedisConfig{
		Addr:     viper.GetString("redis.addr"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       viper.GetInt("redis.db"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize redis: %s", err.Error())
	}

	// 5. Инициализация слоев (Onion Architecture)
	// Репозитории -> Сервисы -> Хендлеры
	repos := repository.NewRepository(db)
	services := service.NewService(repos, redisClient)
	handlers := rest.NewHandler(services, redisClient)

	// 6. Запуск HTTP сервера
	// Заменяем domain.rest на domain.Server
	srv := new(domain.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occurred while running http server: %s", err.Error())
		}
	}()

	logrus.Print("SeeThisGame app started")

	// 7. Graceful Shutdown (Ожидание сигнала завершения)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("SeeThisGame app shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occurred on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occurred on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath(".")       // Поиск в корне
	viper.AddConfigPath("configs") // Поиск в папке configs
	viper.SetConfigName("config")
	viper.SetConfigType("yml")

	// Позволяет переопределять конфиг переменными окружения
	// Например: DB_HOST перекроет значение db.host в yml
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	return viper.ReadInConfig()
}
