package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aintsashqa/go-chat/chat"
	"github.com/aintsashqa/go-chat/internal/controller/http/web"
	"github.com/aintsashqa/go-chat/internal/controller/socket"
	"github.com/aintsashqa/go-chat/internal/repository/redis"
	"github.com/aintsashqa/go-chat/internal/service/render"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const (
	defaultPort    = "8080"
	templatesDir   = "templates"
	templateLayout = "layout"

	// Environment
	envServicePort      = "SERVICE_PORT"
	envLogLevel         = "LOG_LEVEL"
	envRedisURL         = "REDIS_URL"
	envSessionSecretKey = "SESSION_SECRET_KEY"
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

	// Initialze repository
	redisURL := os.Getenv(envRedisURL)
	messageRepository, err := redis.NewMessageRepository(redisURL, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// Initialize services
	cookieService := sessions.NewCookieStore([]byte(os.Getenv(envSessionSecretKey)))
	renderService := render.NewHtmlRenderService(templatesDir, templateLayout, logger)
	messageService := chat.NewMessageService(messageRepository, logger)

	// Initialize chat community
	community := chat.NewCommunity(messageService, logger)

	// Initialize controllers
	socketController := socket.NewSocketController(community, logger)
	hubController := web.NewHubController(community, cookieService, renderService, logger)

	// Initialize routes
	router.Get("/ws", socketController.ReceiveMessage)
	router.Get("/", hubController.NewHub)
	router.Post("/", hubController.CreateHub)
	router.Get("/hub", hubController.JoinHub)
	router.Get("/hub/invite/{hub_id}", hubController.InviteToHub)

	// Starting server on port
	port := servicePort()
	if err := http.ListenAndServe(port, router); err != nil {
		logger.Fatal(err)
	}
}
