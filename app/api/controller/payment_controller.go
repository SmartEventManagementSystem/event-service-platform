package controller

import (
	"net/http"

	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/manager"
	"github.com/labstack/echo/v4"
)

type PaymentController struct {
	res      runtime.Resource
	managers *manager.Managers
}

func NewPaymentController(managers *manager.Managers, res runtime.Resource) *PaymentController {
	return &PaymentController{res: res, managers: managers}
}

func (c *PaymentController) GetPaymentMethods(ec echo.Context) error {
	return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "GetPaymentMethods not implemented"})
}

func (c *PaymentController) CreatePaymentMethod(ec echo.Context) error {
	return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "CreatePaymentMethod not implemented"})
}
