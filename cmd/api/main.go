package main

import (
	"fmt"

	"github.com/andrew-sameh/echo-engine/docs"
	"github.com/andrew-sameh/echo-engine/internal/config"
	"github.com/andrew-sameh/echo-engine/internal/server"
	"github.com/andrew-sameh/echo-engine/internal/server/routes"
	"go.uber.org/zap"
)

//	@title			Echo Engine API
//	@version		1.0
//	@description	This is an Echo API Server template.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath	/api/v1
func main() {
	cfg := config.New()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	server := server.NewServer(cfg)
	routes.RegisterRoutes(server)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	sugar.Infof("Service URL: http://localhost:%s/swagger/index.html", cfg.Server.Port)

	err := server.Start(cfg.Server.Port)

	if err != nil {
		sugar.Fatalf("Cannot start server: %s", err)
	}
}
