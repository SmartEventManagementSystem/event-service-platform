package controller

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type NewsfeedController struct {
    res      runtime.Resource
    managers *manager.Managers
}

func NewNewsfeedController(managers *manager.Managers, res runtime.Resource) *NewsfeedController {
    return &NewsfeedController{res: res, managers: managers}
}

func (c *NewsfeedController) GetUserNewsfeed(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "GetUserNewsfeed not implemented"})
}

func (c *NewsfeedController) InvalidateNewsfeed(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "InvalidateNewsfeed not implemented"})
}
