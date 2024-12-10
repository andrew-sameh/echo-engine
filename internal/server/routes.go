package server

import (
	"net/http"

	_ "github.com/andrew-sameh/echo-engine/docs"
	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
)

//	@title			Echo Engine API
//	@version		1.0
//	@description	This is a sample server Petstore server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/api/v1
func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.RequestID())
	zapLogger, _ := zap.NewProduction()
	e.Use(echozap.ZapLogger(zapLogger))
	// e.Use(middleware.Timeout())

	// e.Use(middleware.Logger())

	e.Use(middleware.Recover())
	e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		zapLogger.Info("Request Body", zap.String("body", string(reqBody)), zap.String("path", c.Path()), zap.String("method", c.Request().Method), zap.String("query", c.QueryString()), zap.String("remote_ip", c.RealIP()), zap.String("host", c.Request().Host), zap.String("user_agent", c.Request().UserAgent()), zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
		zapLogger.Info("Response Body", zap.String("body", string(resBody)), zap.String("path", c.Path()), zap.String("method", c.Request().Method), zap.String("query", c.QueryString()), zap.String("remote_ip", c.RealIP()), zap.String("host", c.Request().Host), zap.String("user_agent", c.Request().UserAgent()), zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)))
	}))
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	base := e.Group("/api/v1")

	base.GET("/", s.HelloWorldHandler)

	base.GET("/health", s.healthHandler)

	base.GET("/users", s.ListUsersHandler)
	zapLogger.Info("Routes registered")

	return e
}

// HelloWorldHandler returns a Hello World message
//
//	@Summary		Hello World
//	@Description	Returns a Hello World message
//	@Tags			HelloWorld
//	@Accept			json
//	@Produce		json
//	@Success		200	{object} map[string]string
//	@Router			/ [get]
func (s *Server) HelloWorldHandler(c echo.Context) error {
	logger := c.Logger()
	resp := map[string]string{
		"message": "Hello World",
	}
	logger.Infof("Hello World")
	return c.JSON(http.StatusOK, resp)
}

// ListUsersHandler lists all existing users
//
//	@Summary		List users
//	@Description	get users
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array} []map[string]string
//	@Failure		500	{object} error
//	@Router			/users [get]
func (s *Server) ListUsersHandler(c echo.Context) error {
	queries := s.db.Queries()
	users, err := queries.GetAllUsers(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, users)
}

// healthHandler checks the health of the server
//
//	@Summary		Health check
//	@Description	Checks the health of the server
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
