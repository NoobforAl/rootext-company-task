package router

import (
	"ratblog/config"
	"ratblog/internal/services/v1/controller"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func setupPostsRouter(e *echo.Group, config config.Config, con *controller.Controller) {
	g := e.Group("/posts")

	g.GET("/", con.GetPosts())
	g.GET("/filter", con.GetPostsByFilter())

	g.GET("/:id", con.GetPost())
	g.POST("/", con.CreatePost(), echojwt.JWT([]byte(config.JWT.Secret)))
	g.PUT("/:id", con.Update(), echojwt.JWT([]byte(config.JWT.Secret)))
	g.DELETE("/:id", con.DeletePost(), echojwt.JWT([]byte(config.JWT.Secret)))

	g.GET("/:id/up", con.UpVotePost(), echojwt.JWT([]byte(config.JWT.Secret)))
	g.GET("/:id/down", con.DownVotePost(), echojwt.JWT([]byte(config.JWT.Secret)))
	g.GET("/:id/removeRate", con.RemoveRatePost(), echojwt.JWT([]byte(config.JWT.Secret)))
}
