package api

import (
	"fmt"
	"net/http"

	"github.com/Francesco99975/trapk/internal/models"
	"github.com/labstack/echo/v4"
)

func Report() echo.HandlerFunc {
	return func(c echo.Context) error {
		var payload models.Report
		if err := c.Bind(&payload); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error parsing data for report: %v", err), Errors: []string{err.Error()}})
		}
		ip := c.RealIP()

		if err := models.CreateReport(payload, ip); err != nil {
			return c.JSON(http.StatusBadRequest, models.JSONErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("Error creating report: %v", err), Errors: []string{err.Error()}})
		}

		return c.JSON(http.StatusCreated, "OK")
	}
}
