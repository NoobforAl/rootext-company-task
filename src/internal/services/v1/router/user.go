package router

import "github.com/labstack/echo/v4"

func setupUsersRouter(g *echo.Group) {
	g.GET("/users/", emttp)
	g.POST("/users/", emttp)
	g.GET("/users/:id", emttp)
	g.PUT("/users/:id", emttp)
	g.DELETE("/users/:id", emttp)
}
