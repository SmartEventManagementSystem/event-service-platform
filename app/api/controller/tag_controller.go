package controller

import (
	"net/http"

	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/manager"
	"github.com/labstack/echo/v4"
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
