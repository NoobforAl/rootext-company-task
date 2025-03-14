package router

import "github.com/labstack/echo/v4"

var ttt echo.HandlerFunc

func SetupRouter(e *echo.Echo) {
	eg := e.Group("/v1")

	setupUsersRouter(eg)
	setupPostsRouter(eg)

	eg.GET("/health", ttt)
}
