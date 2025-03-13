package router

import "github.com/labstack/echo/v4"

var emttp echo.HandlerFunc

func setupPostsRouter(g *echo.Group) {
	g.GET("/posts/", emttp)
	g.POST("/posts/", emttp)
	g.GET("/posts/:id", emttp)
	g.PUT("/posts/:id", emttp)
	g.DELETE("/posts/:id", emttp)

	g.GET("/posts/:id/up", emttp)
	g.GET("/posts/:id/down", emttp)

	g.GET("/posts/:id/comments/", emttp)
	g.POST("/posts/:id/comments/", emttp)

	g.GET("/posts/:id/comments/:comment_id", emttp)
	g.PUT("/posts/:id/comments/:comment_id", emttp)
	g.DELETE("/posts/:id/comments/:comment_id", emttp)
}
