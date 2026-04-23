package app

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vinzryyy/iot-backend/models"
	"github.com/vinzryyy/iot-backend/service"
)

type AuthHandler struct {
	auth *service.AuthService
}

func NewAuthHandler(a *service.AuthService) *AuthHandler {
	return &AuthHandler{auth: a}
}

// Login godoc
// @Summary      Login
// @Description  Authenticate with email and password, returns a JWT
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "credentials"
// @Success      200 {object} models.APIResponse{data=models.LoginResponse}
// @Failure      400 {object} models.APIResponse
// @Failure      401 {object} models.APIResponse
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	resp, err := h.auth.Login(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "login successful",
		Data:    resp,
	})
}

// Register godoc
// @Summary      Register user (admin only)
// @Description  Create a new user with optional location access. Requires admin role.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.RegisterRequest true "new user"
// @Security     BearerAuth
// @Success      201 {object} models.APIResponse{data=models.User}
// @Failure      400 {object} models.APIResponse
// @Failure      409 {object} models.APIResponse
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req models.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	u, err := h.auth.Register(c.Request().Context(), req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "user created",
		Data:    u,
	})
}

// Me godoc
// @Summary      Get current user profile
// @Tags         user
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.APIResponse{data=models.UserProfile}
// @Failure      401 {object} models.APIResponse
// @Router       /me [get]
func (h *AuthHandler) Me(c echo.Context) error {
	userID, _ := c.Get(CtxUserID).(string)
	profile, err := h.auth.Profile(c.Request().Context(), userID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    profile,
	})
}
