package main

import (
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/mqqff/savebite-be/internal/domain/env"
	"github.com/mqqff/savebite-be/internal/infra/mysql"
	"github.com/mqqff/savebite-be/internal/infra/server"
)

func main() {
	db, err := mysql.NewConn()
	if err != nil {
		return
	}

	httpServer := server.NewHTTPServer()
	app := httpServer.GetApp()

	httpServer.MountMiddlewares()

	app.Get("/metrics", monitor.New())

	httpServer.MountRoutes(db)
	httpServer.Start(env.AppEnv.AppHost + ":" + env.AppEnv.AppPort)
}
