package main

import (
	"io"

	"github.com/gin-gonic/gin"

	"fliqt/config"
	"fliqt/internal/handler"
	"fliqt/internal/util"
)

func main() {
	cfg := config.NewConfig()
	logger := util.NewLogger(cfg)

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
	}
	app := gin.Default()

	db, err := util.NewGormDB(cfg)
	if err != nil {
		panic(err)
	}

	app.Use(handler.Logger(logger))
	app.Use(handler.ErrorHandler(logger))
	app.NoRoute(handler.NotFoundHandler())

	handler.NewRouter(
		app,
		logger,
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
