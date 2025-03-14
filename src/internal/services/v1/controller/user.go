package controller

import (
	"ratblog/config"
	"ratblog/internal/entity"
	"ratblog/internal/services/v1/schema"

	"github.com/labstack/echo/v4"
)

type jwtToken struct {
	Token string `json:"token"`
}

func (con *Controller) Login(config config.Config) echo.HandlerFunc {
	return func(c echo.Context) error {
		var loginInfo schema.Login
		if err := c.Bind(&loginInfo); err != nil {
			return c.JSON(400, badRequestResponse)
		}

		err := loginInfo.Validate()
		if err != nil {
			return c.JSON(400, response{Message: err.Error()})
		}

		ctx := c.Request().Context()
		user, err := con.repo.CheckUserExistWithUsernamePassword(ctx, loginInfo.UserName, loginInfo.Password)
		if err != nil {
			return c.JSON(500, usernameOrPasswordIncorrectResponse)
		}

		token, err := generateToken(int(user.ID), user.Username, user.Role, []byte(config.JWT.Secret), config.JWT.Expiration)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		return c.JSON(200, jwtToken{Token: token})
	}
}

func (con *Controller) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		var registerUser schema.Register
		if err := c.Bind(&registerUser); err != nil {
			return c.JSON(400, badRequestResponse)
		}

		if err := registerUser.Validate(); err != nil {
			return c.JSON(400, response{Message: err.Error()})
		}

		newUser := &entity.User{
			Email:     registerUser.Email,
			Username:  registerUser.UserName,
			FirstName: registerUser.FirstName,
			LastName:  registerUser.LastName,
			Password:  registerUser.Password,
			Role:      "user",
		}

		ctx := c.Request().Context()
		err := con.repo.CreateUser(ctx, newUser)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		return c.JSON(200, response{Message: "User created successfully"})
	}
}

func (con *Controller) GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := validateUser(c)
		if err != nil {
			return c.JSON(401, unauthorizedResponse)
		}

		ctx := c.Request().Context()
		userData, err := con.repo.GetUserByID(ctx, user.ID)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		userInfo := schema.UserInfo{
			Email:     userData.Email,
			UserName:  userData.Username,
			FirstName: userData.FirstName,
			LastName:  userData.LastName,
			Role:      user.Role,
		}

		return c.JSON(200, userInfo)
	}
}

func (con *Controller) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := validateUser(c)
		if err != nil {
			return c.JSON(401, unauthorizedResponse)
		}

		var updateInfo schema.UserUpdate
		if err := c.Bind(&updateInfo); err != nil {
			return c.JSON(400, badRequestResponse)
		}

		userEntity := &entity.User{
			ID:        int32(user.ID),
			FirstName: updateInfo.FirstName,
			LastName:  updateInfo.LastName,
			Password:  updateInfo.Password,
		}

		ctx := c.Request().Context()
		err = con.repo.UpdateUser(ctx, userEntity)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		return c.JSON(200, response{Message: "User updated successfully"})
	}
}

func (con *Controller) Delete() echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := validateUser(c)
		if err != nil {
			return c.JSON(401, unauthorizedResponse)
		}

		ctx := c.Request().Context()
		err = con.repo.DeleteUser(ctx, user.ID)
		if err != nil {
			return c.JSON(500, response{Message: err.Error()})
		}

		return c.JSON(200, response{Message: "User deleted successfully"})
	}
}

func (con *Controller) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(400, notImplementedResponse)
	}
}

func (con *Controller) GetUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(400, notImplementedResponse)
	}
}

func (con *Controller) UpdateUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(400, notImplementedResponse)
	}
}

func (con *Controller) DeleteUserById() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(400, notImplementedResponse)
	}
}
