package app

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/vinzryyy/iot-backend/service"
)

const (
	CtxUserID = "uid"
	CtxEmail  = "email"
	CtxRole   = "role"
)

func JWTAuth(j *service.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid Authorization header")
			}
			tokenStr := strings.TrimSpace(auth[len("Bearer "):])
			claims, err := j.Parse(tokenStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
			}
			c.Set(CtxUserID, claims.UserID)
			c.Set(CtxEmail, claims.Email)
			c.Set(CtxRole, claims.Role)
			return next(c)
		}
	}
}

// RequireRole blocks requests whose role is not in the allow list.
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			role, _ := c.Get(CtxRole).(string)
			for _, r := range roles {
				if r == role {
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusForbidden, "insufficient role")
		}
	}
}
