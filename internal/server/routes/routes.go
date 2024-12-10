package routes

import (
	"time"

	_ "github.com/andrew-sameh/echo-engine/docs"
	s "github.com/andrew-sameh/echo-engine/internal/server"
	h "github.com/andrew-sameh/echo-engine/internal/server/handlers"
	"github.com/brpaz/echozap"

	// "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

func RegisterRoutes(s *s.Server) {
	zapLogger, _ := zap.NewProduction()

	// Server config
	s.Echo.Server.ReadTimeout = 10 * time.Second
	s.Echo.Server.WriteTimeout = 30 * time.Second
	s.Echo.Server.IdleTimeout = time.Minute

	// Handlers creation
	genericHandler := h.NewGenericHandler(s)

	// Middlewares
	s.Echo.Use(middleware.RequestID())
	s.Echo.Use(echozap.ZapLogger(zapLogger))
	s.Echo.Use(middleware.CORS())
	s.Echo.Use(middleware.Recover())
	// e.Use(middleware.Timeout())
	// s.Echo.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	zapLogger.Info("Request Body", zap.String("body", string(reqBody)), zap.String("path", c.Path()), zap.String("method", c.Request().Method), zap.String("query", c.QueryString()), zap.String("remote_ip", c.RealIP()), zap.String("host", c.Request().Host), zap.String("user_agent", c.Request().UserAgent()), zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
	// 	zapLogger.Info("Response Body", zap.String("body", string(resBody)), zap.String("path", c.Path()), zap.String("method", c.Request().Method), zap.String("query", c.QueryString()), zap.String("remote_ip", c.RealIP()), zap.String("host", c.Request().Host), zap.String("user_agent", c.Request().UserAgent()), zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
	// }))

	// Base Routes
	s.Echo.GET("/swagger/*", echoSwagger.WrapHandler)

	// Versioned Routes
	base := s.Echo.Group("/api/v1")

	base.GET("/", genericHandler.HelloWorldHandler)
	base.GET("/health", genericHandler.HealthHandler)
	base.GET("/users", genericHandler.ListUsersHandler)

}
