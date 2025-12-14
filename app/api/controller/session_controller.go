package controller

import (
	"net/http"

	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/manager"
	"github.com/labstack/echo/v4"
)

type SessionController struct {
	res      runtime.Resource
	managers *manager.Managers
}

func NewSessionController(managers *manager.Managers, res runtime.Resource) *SessionController {
	return &SessionController{res: res, managers: managers}
}

func (c *SessionController) CreateSession(ec echo.Context) error {
	return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "CreateSession not implemented"})
}

func (c *SessionController) ListSessions(ec echo.Context) error {
	return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "ListSessions not implemented"})
}
