package controller

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type SubscriptionController struct {
    res      runtime.Resource
    managers *manager.Managers
}

func NewSubscriptionController(managers *manager.Managers, res runtime.Resource) *SubscriptionController {
    return &SubscriptionController{res: res, managers: managers}
}

func (c *SubscriptionController) CreateSubscription(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "CreateSubscription not implemented"})
}

func (c *SubscriptionController) GetUserSubscription(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "GetUserSubscription not implemented"})
}
