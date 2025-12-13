package controller

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type GroupController struct {
    res      runtime.Resource
    managers *manager.Managers
}

func NewGroupController(managers *manager.Managers, res runtime.Resource) *GroupController {
    return &GroupController{res: res, managers: managers}
}

func (c *GroupController) CreateGroup(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "CreateGroup not implemented"})
}

func (c *GroupController) GetGroups(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "GetGroups not implemented"})
}
