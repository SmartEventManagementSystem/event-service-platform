package controller

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type ChatController struct {
    res      runtime.Resource
    managers *manager.Managers
}

func NewChatController(managers *manager.Managers, res runtime.Resource) *ChatController {
    return &ChatController{res: res, managers: managers}
}

func (c *ChatController) CreateRoom(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "CreateRoom not implemented"})
}

func (c *ChatController) ListRooms(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "ListRooms not implemented"})
}

func (c *ChatController) SendMessage(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "SendMessage not implemented"})
}
