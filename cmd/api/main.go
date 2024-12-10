package main

import (
	"os"

	"github.com/andrew-sameh/echo-engine/internal/server"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	server := server.NewServer()
	sugar.Infof("Starting server on port %s", os.Getenv("PORT"))
	// Log service url host and port
	sugar.Infof("Service URL: http://localhost:%s/swagger/index.html", os.Getenv("port"))

	err := server.ListenAndServe()
	if err != nil {
		// panic(fmt.Sprintf("cannot start server: %s", err))
		sugar.Errorf("cannot start server: %s", err)
	}
}
