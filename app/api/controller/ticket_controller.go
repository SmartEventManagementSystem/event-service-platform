package controller

import (
    "net/http"

    "github.com/labstack/echo/v4"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/internal/runtime"
    "github.com/iamhuutho/Event-Management-System/Event-Management-System-BE/event-service-platform/app/manager"
)

type TicketController struct {
    res      runtime.Resource
    managers *manager.Managers
}

func NewTicketController(managers *manager.Managers, res runtime.Resource) *TicketController {
    return &TicketController{res: res, managers: managers}
}

func (c *TicketController) CreateTicketSale(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "CreateTicketSale not implemented"})
}

func (c *TicketController) ListTickets(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "ListTickets not implemented"})
}

func (c *TicketController) ValidateTicket(ec echo.Context) error {
    return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "ValidateTicket not implemented"})
}
