package alert

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) TriggerPriceAlert(c *gin.Context) {

	alerts, err := h.s.GetTokenAlerts(c)
	if err != nil {
		c.JSON(500, Logger{
			Message: err.Error(),
		})
		return
	}

	err = h.s.SendAlerts(c, alerts)
	if err != nil {
		c.JSON(500, Logger{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, "OK")
}
