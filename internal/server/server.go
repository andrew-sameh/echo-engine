package server

import (
	"github.com/andrew-sameh/echo-engine/internal/config"
	db "github.com/andrew-sameh/echo-engine/internal/database"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo   *echo.Echo
	Config *config.Config

	DB db.DBService
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Config: cfg,
		Echo:   echo.New(),
		DB:     db.NewConnection(),
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
