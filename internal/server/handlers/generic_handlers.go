package handlers

import (
	"net/http"

	s "github.com/andrew-sameh/echo-engine/internal/server"
	"github.com/labstack/echo/v4"
)

type GenericHandler struct {
	server *s.Server
}

func NewGenericHandler(server *s.Server) *GenericHandler {
	return &GenericHandler{server: server}
}

// HelloWorldHandler returns a Hello World message
//
//	@Summary		Hello World
//	@Description	Returns a Hello World message
//	@Tags			Generic
//	@Accept			json
//	@Produce		json
//	@Success		200	{object} map[string]string
//	@Router			/ [get]
func (g *GenericHandler) HelloWorldHandler(c echo.Context) error {
	logger := c.Logger()
	resp := map[string]string{
		"message": "Hello World",
	}
	logger.Infof("Hello World")
	return c.JSON(http.StatusOK, resp)
}

// healthHandler checks the health of the server
//
//	@Summary		Health check
//	@Description	Checks the health of the server
//	@Tags			Generic
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (g *GenericHandler) HealthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, g.server.DB.Health())
}
