package api

import (
	"net/http"

	"github.com/Francesco99975/trapk/internal/models"
	"github.com/labstack/echo/v4"
)

func Report() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.Report
		if err := c.Bind(&payload); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
		}
		ip := c.RealIP()

		if err := models.CreateReport(payload, ip); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Could not create report")
		}

		return c.JSON(http.StatusCreated, "OK")
	}
}
