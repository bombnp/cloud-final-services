package summary

import (
	"log"

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
		c.JSON(400, &Logger{
			Message: err.Error(),
		})
		return
	}

	for address, summary := range summaryMap {
		log.Println(address, summary)
	}

	c.JSON(200, &Logger{
		Message: "OK",
	})

}
