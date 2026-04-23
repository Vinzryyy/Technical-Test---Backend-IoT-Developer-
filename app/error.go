package app

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/vinzryyy/iot-backend/models"
	"github.com/vinzryyy/iot-backend/repo"
	"github.com/vinzryyy/iot-backend/service"
)

// ErrorHandler translates application errors into standard JSON envelopes.
func ErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	status := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, repo.ErrNotFound):
		status, msg = http.StatusNotFound, "resource not found"
	case errors.Is(err, service.ErrInvalidCredentials):
		status, msg = http.StatusUnauthorized, err.Error()
	case errors.Is(err, service.ErrForbidden):
		status, msg = http.StatusForbidden, err.Error()
	case errors.Is(err, service.ErrEmailExists):
		status, msg = http.StatusConflict, err.Error()
	}

	var he *echo.HTTPError
	if errors.As(err, &he) {
		status = he.Code
		if m, ok := he.Message.(string); ok {
			msg = m
		} else {
			msg = http.StatusText(status)
		}
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		status = http.StatusBadRequest
		msg = "validation failed: " + ve.Error()
	}

	_ = c.JSON(status, models.APIResponse{
		Success: false,
		Error:   msg,
	})
}
