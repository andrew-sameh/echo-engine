package routes

import (
	"time"

	_ "github.com/andrew-sameh/echo-engine/docs"
	s "github.com/andrew-sameh/echo-engine/internal/server"
	h "github.com/andrew-sameh/echo-engine/internal/server/handlers"
	"github.com/andrew-sameh/echo-engine/internal/services/token"
	"github.com/andrew-sameh/echo-engine/pkg/logger"
	"github.com/brpaz/echozap"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func RegisterRoutes(s *s.Server) {
	// Server config
	s.Echo.Server.ReadTimeout = 10 * time.Second
	s.Echo.Server.WriteTimeout = 30 * time.Second
	s.Echo.Server.IdleTimeout = time.Minute
	s.Echo.Validator = s.Config.Server.Validator
	s.Echo.Binder = s.Config.Server.Binder

	// Handlers creation
	genericHandler := h.NewGenericHandler(s)
	authHandler := h.NewAuthHandler(s)
	userHandler := h.NewUserHandler(s)

	// Middlewares
	s.Echo.Use(middleware.RequestID())
	s.Echo.Use(echozap.ZapLogger(logger.ZapLogger()))
	s.Echo.Use(middleware.CORSWithConfig(s.Config.Server.CORSConfig))
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

	auth := base.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)
	auth.POST("/refresh", authHandler.RefreshToken)

	r := base.Group("")
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		NewClaimsFunc: func(_ echo.Context) jwt.Claims {
			return new(token.JwtCustomClaims)
		},
		SigningKey: []byte(s.Config.Auth.AccessSecret),
	}
	r.Use(echojwt.WithConfig(config))
	r.GET("/profile", userHandler.GetMyUserHandler)
	r.GET("/users", userHandler.ListUsersHandler)

}
