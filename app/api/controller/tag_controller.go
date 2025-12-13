package controller

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type TagController struct {
    res      runtime.Resource
    managers *manager.Managers
}

func NewTagController(managers *manager.Managers, res runtime.Resource) *TagController {
    return &TagController{res: res, managers: managers}
}

func (c *TagController) GetAllTags(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "GetAllTags not implemented"})
}

func (c *TagController) CreateTag(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "CreateTag not implemented"})
}
