package controller

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type WorkspaceController struct {
    res      runtime.Resource
    managers *manager.Managers
}

func NewWorkspaceController(managers *manager.Managers, res runtime.Resource) *WorkspaceController {
    return &WorkspaceController{res: res, managers: managers}
}

func (c *WorkspaceController) CreateWorkspace(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "CreateWorkspace not implemented"})
}

func (c *WorkspaceController) ListWorkspaces(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "ListWorkspaces not implemented"})
}
