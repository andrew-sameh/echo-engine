package server

import (
	"net/http"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	zapLogger, _ := zap.NewProduction()
	e.Use(middleware.RequestID())
	// e.Use(middleware.Timeout())

	// e.Use(middleware.Logger())
	e.Use(echozap.ZapLogger(zapLogger))

	e.Use(middleware.Recover())

	e.GET("/", s.HelloWorldHandler)

	e.GET("/health", s.healthHandler)

	e.GET("/users", s.ListUsersHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) ListUsersHandler(c echo.Context) error {
	queries := s.db.Queries()
	users, err := queries.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
