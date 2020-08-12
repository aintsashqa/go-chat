package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aintsashqa/go-chat/chat"
	"github.com/aintsashqa/go-chat/internal/api"
	"github.com/aintsashqa/go-chat/internal/repository/redis"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort  = "8080"
	templatesDir = "templates"

	// Environment
	envServicePort = "SERVICE_PORT"
	envLogLevel    = "LOG_LEVEL"
	envRedisURL    = "REDIS_URL"
)

func servicePort() string {
	port := defaultPort
	if os.Getenv(envServicePort) != "" {
		port = os.Getenv(envServicePort)
	}
	return fmt.Sprintf(":%s", port)
}

func logLevel() logrus.Level {
	level, err := logrus.ParseLevel(os.Getenv(envLogLevel))
	if err != nil {
		return logrus.InfoLevel
	}

	return level
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func main() {
	logger := logrus.New()
	logger.Level = logLevel()

	router := chi.NewRouter()
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	redisURL := os.Getenv(envRedisURL)
	messageRepository, err := redis.NewMessageRepository(redisURL, logger)
	if err != nil {
		logger.Fatal(err)
	}

	messageService := chat.NewMessageService(messageRepository, logger)
	hub := chat.NewHub(messageService, logger)
	go hub.Run()

	wsHandler := api.NewWSHandler(hub, logger)
	httpHandler := api.NewHttpHandler(templatesDir, logger)

	router.Get("/", httpHandler.JoinChat)
	router.Get("/ws", wsHandler.ReceiveMessage)

	port := servicePort()
	if err := http.ListenAndServe(port, router); err != nil {
		logger.Fatal(err)
	}
}
