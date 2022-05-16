package summary

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

func (h *Handler) TriggerDailySummaryReport(c *gin.Context) {

	summaryMap, err := h.s.GetTokenDailySummary(c)
	if err != nil {
		c.JSON(500, &Logger{
			Message: err.Error(),
		})
		return
	}

	err = h.s.SendSummaryReports(c, summaryMap)
	if err != nil {
		c.JSON(500, &Logger{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, &Logger{
		Message: "OK",
	})

}
