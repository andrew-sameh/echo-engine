package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	db "github.com/andrew-sameh/echo-engine/init"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int

	db db.DBService
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db: db.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
