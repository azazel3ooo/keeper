package server_repo

import (
	"github.com/azazel3ooo/keeper/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	"net/http"
)

// SetupApp производит настройку fiber.App и инициализирует хендлеры
func (s *Server) SetupApp() {
	a := fiber.New()

	a.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed,
	}))
	a.Use(recover.New(recover.Config{EnableStackTrace: true}))
	a.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path} ${resBody}\n",
	}))

	api := a.Group("/api")
	v1 := api.Group("/v1")

	v1.Get("/items", s.getAll)
	v1.Post("/items", s.set)
	v1.Delete("/items", s.delete)
	v1.Patch("/items", s.update)

	v1.Post("/registration", s.registration)
	v1.Post("/auth", s.authorization)

	v1.Get("/swagger/*", fiberSwagger.WrapHandler)

	s.app = a
}

// Listen запускает приложение на указанному порту
func (s Server) Listen() error {
	return s.app.Listen(s.cfg.HostAddr)
}

// NewServer возвращает сервер, применяя к нему указанные опции
func NewServer(opts ...func(*Server)) *Server {
	s := &Server{}
	s.app = fiber.New()

	for _, opt := range opts {
		opt(s)
	}
	return s
}

// WithProcessingChan добавляет ProcessingChan серверу
func WithProcessingChan(c ProcessingChan) func(*Server) {
	return func(s *Server) {
		s.processingChan = c
	}
}

func WithLoggingChan(c LogChan) func(*Server) {
	return func(s *Server) {
		s.lgChan = c
	}
}

func WithLogger(l any) func(*Server) {
	return func(s *Server) {
		s.lg = l
	}
}

// WithStorage добавляет Storable4Server серверу
func WithStorage(store models.Storable4Server) func(*Server) {
	return func(s *Server) {
		s.storage = store
	}
}

// WithConfig добавляет Config серверу
func WithConfig(c models.Config) func(*Server) {
	return func(s *Server) {
		s.cfg = c
	}
}

// Test производит тестовый запрос к серверу, используя инструмент фреймворка. (нужно для тестирования клиента)
func (s Server) Test(req *http.Request, timeout int) (*http.Response, error) {
	return s.app.Test(req, timeout)
}
