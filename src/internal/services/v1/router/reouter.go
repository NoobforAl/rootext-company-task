package router

import (
	"ratblog/config"
	"ratblog/contract"
	"ratblog/internal/services/v1/controller"

	"github.com/labstack/echo/v4"
)

func SetupRouter(e *echo.Echo, repo contract.Repository, config config.Config) {
	con := controller.New(repo, config)

	eg := e.Group("/api/v1")

	// this endpoint for see tags for just tests
	eg.GET("/tags", con.GetTags())

	setupUsersRouter(eg, config, con)
	setupPostsRouter(eg, config, con)

	eg.GET("/health", func(c echo.Context) error {
		return c.JSON(200, "OK")
	})
}
