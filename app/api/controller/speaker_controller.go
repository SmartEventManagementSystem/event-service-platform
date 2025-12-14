package controller

import (
	"net/http"

	"backend/event-service-platform/app/internal/runtime"
	"backend/event-service-platform/app/manager"
	"github.com/labstack/echo/v4"
)

type SpeakerController struct {
	res      runtime.Resource
	managers *manager.Managers
}

func NewSpeakerController(managers *manager.Managers, res runtime.Resource) *SpeakerController {
	return &SpeakerController{res: res, managers: managers}
}

func (c *SpeakerController) CreateSpeaker(ec echo.Context) error {
	return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "CreateSpeaker not implemented"})
}

func (c *SpeakerController) ListSpeakers(ec echo.Context) error {
	return ec.JSON(http.StatusNotImplemented, map[string]string{"error": "ListSpeakers not implemented"})
}
