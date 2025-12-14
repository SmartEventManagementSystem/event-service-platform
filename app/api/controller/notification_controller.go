package controller

import (
	"net/http"

	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/manager"
	"github.com/labstack/echo/v4"
)

type NotificationController struct {
	res      runtime.Resource
	managers *manager.Managers
}

func NewNotificationController(managers *manager.Managers, res runtime.Resource) *NotificationController {
	return &NotificationController{res: res, managers: managers}
}

func (c *NotificationController) GetNotifications(ec echo.Context) error {
	return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "GetNotifications not implemented"})
}

func (c *NotificationController) MarkAsRead(ec echo.Context) error {
	return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "MarkAsRead not implemented"})
}
