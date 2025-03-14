package controller

import (
	"net/http"
	"ratblog/internal/services/v1/schema"

	"github.com/labstack/echo/v4"
)

func (con *Controller) GetTags() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		raw, err := con.repo.GetAllTags(ctx)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		var tags []*schema.Tag
		for _, r := range raw {
			tags = append(tags, &schema.Tag{
				ID:        int(r.ID),
				Name:      r.Name,
				CreatedAt: r.CreatedAt,
				UpdatedAt: r.UpdatedAt,
			})
		}

		return c.JSON(http.StatusOK, tags)
	}
}
