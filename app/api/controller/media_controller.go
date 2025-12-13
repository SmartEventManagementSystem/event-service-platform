package controller

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type MediaController struct {
    res      runtime.Resource
    managers *manager.Managers
}

func NewMediaController(managers *manager.Managers, res runtime.Resource) *MediaController {
    return &MediaController{res: res, managers: managers}
}

func (c *MediaController) UploadChatFile(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "UploadChatFile not implemented"})
}

func (c *MediaController) GetChatAttachments(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "GetChatAttachments not implemented"})
}
