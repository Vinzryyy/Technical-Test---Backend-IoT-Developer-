package app

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/vinzryyy/iot-backend/models"
	"github.com/vinzryyy/iot-backend/service"
)

type DeviceHandler struct {
	devices *service.DeviceService
}

func NewDeviceHandler(d *service.DeviceService) *DeviceHandler {
	return &DeviceHandler{devices: d}
}

func (h *DeviceHandler) actor(c echo.Context) (userID, role string) {
	userID, _ = c.Get(CtxUserID).(string)
	role, _ = c.Get(CtxRole).(string)
	return
}

func parseIntID(c echo.Context) (int64, error) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, echo.NewHTTPError(http.StatusBadRequest, "id must be a positive integer")
	}
	return id, nil
}

// List godoc
// @Summary      List devices
// @Description  Returns only devices in the caller's accessible locations
// @Tags         devices
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.APIResponse{data=[]models.Device}
// @Router       /devices [get]
func (h *DeviceHandler) List(c echo.Context) error {
	uid, role := h.actor(c)
	items, err := h.devices.List(c.Request().Context(), uid, role)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: items})
}

// Get godoc
// @Summary      Get device by ID
// @Tags         devices
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Device ID"
// @Success      200 {object} models.APIResponse{data=models.Device}
// @Failure      403 {object} models.APIResponse
// @Failure      404 {object} models.APIResponse
// @Router       /devices/{id} [get]
func (h *DeviceHandler) Get(c echo.Context) error {
	id, err := parseIntID(c)
	if err != nil {
		return err
	}
	uid, role := h.actor(c)
	d, err := h.devices.Get(c.Request().Context(), uid, role, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.APIResponse{Success: true, Data: d})
}

// Create godoc
// @Summary      Create device
// @Tags         devices
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.CreateDeviceRequest true "device"
// @Success      201 {object} models.APIResponse{data=models.Device}
// @Failure      400 {object} models.APIResponse
// @Failure      403 {object} models.APIResponse
// @Router       /devices [post]
func (h *DeviceHandler) Create(c echo.Context) error {
	var req models.CreateDeviceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	uid, role := h.actor(c)
	d, err := h.devices.Create(c.Request().Context(), uid, role, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, models.APIResponse{
		Success: true, Message: "device created", Data: d,
	})
}

// Update godoc
// @Summary      Update device
// @Tags         devices
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Device ID"
// @Param        request body models.UpdateDeviceRequest true "device"
// @Success      200 {object} models.APIResponse{data=models.Device}
// @Router       /devices/{id} [put]
func (h *DeviceHandler) Update(c echo.Context) error {
	id, err := parseIntID(c)
	if err != nil {
		return err
	}
	var req models.UpdateDeviceRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}
	uid, role := h.actor(c)
	d, err := h.devices.Update(c.Request().Context(), uid, role, id, req)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true, Message: "device updated", Data: d,
	})
}

// Delete godoc
// @Summary      Delete device
// @Tags         devices
// @Security     BearerAuth
// @Param        id path int true "Device ID"
// @Success      200 {object} models.APIResponse
// @Router       /devices/{id} [delete]
func (h *DeviceHandler) Delete(c echo.Context) error {
	id, err := parseIntID(c)
	if err != nil {
		return err
	}
	uid, role := h.actor(c)
	if err := h.devices.Delete(c.Request().Context(), uid, role, id); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, models.APIResponse{
		Success: true, Message: "device deleted",
	})
}
