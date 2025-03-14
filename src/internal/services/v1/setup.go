package serviceV1

import (
	"context"
	"ratblog/config"
	"ratblog/contract"
	"ratblog/internal/services/v1/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupHttpServerV1(ctx context.Context, repo contract.Repository, cfg config.Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	router.SetupRouter(e, repo, cfg)
	return e
}
