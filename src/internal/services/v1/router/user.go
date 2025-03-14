package router

import (
	"ratblog/config"
	"ratblog/internal/services/v1/controller"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func setupUsersRouter(e *echo.Group, config config.Config, con *controller.Controller) {
	g := e.Group("/users")

	g.POST("/login", con.Login(config))
	g.POST("/register", con.Register())

	g.GET("/userInfo", con.GetUser(), echojwt.JWT([]byte(config.JWT.Secret)))
	g.PUT("/", con.Update(), echojwt.JWT([]byte(config.JWT.Secret)))
	g.DELETE("/", con.Delete(), echojwt.JWT([]byte(config.JWT.Secret)))

	// TODO: admin permission
	g.GET("/", con.GetUsers())
	g.GET("/:id", con.GetUserById())
	g.PUT("/:id", con.UpdateUserById())
	g.DELETE("/:id", con.DeleteUserById())
}
