package main

import (
	"fmt"

	"github.com/andrew-sameh/echo-engine/docs"
	"github.com/andrew-sameh/echo-engine/internal/config"
	"github.com/andrew-sameh/echo-engine/internal/server"
	"github.com/andrew-sameh/echo-engine/internal/server/routes"
	"github.com/andrew-sameh/echo-engine/pkg/logger"
)

//	@title			Echo Engine API
//	@version		1.0
//	@description	This is an Echo API Server template.

//	@contact.name	Andrew Sameh
//	@contact.url	https://andrewsam.xyz
//	@contact.email	g.andrewsameh@gmail.com

//	@license.name	MIT
//	@license.url	https://opensource.org/license/mit

//	@securityDefinitions.apiKey ApiKeyAuth
//	@in							header
//	@name						Authorization

// @BasePath	/api/v1
func main() {
	cfg := config.New()

	zlog := logger.NewLogger(cfg.Logger)
	if zlog.Zap != nil {
		defer zlog.Zap.Sync()
	}

	server := server.NewServer(cfg, zlog)
	routes.RegisterRoutes(server)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	zlog.Zap.Infof("Service URL: http://localhost:%s/swagger/index.html", cfg.Server.Port)

	err := server.Start(cfg.Server.Port)

	if err != nil {
		zlog.Zap.Fatalf("Cannot start server: %s", err)
	}
}
