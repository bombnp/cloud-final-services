package alert

import "github.com/gin-gonic/gin"

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) GetTokenAlertSummaryHandler(c *gin.Context) {

	summary_map, err := h.s.GetTokenAlertSummary(c)

	if err != nil {
		c.JSON(400, Logger{
			Message: err.Error(),
		})
		return
	}

	var response []AlertResponse

	for address, summary := range summary_map {
		response = append(response, AlertResponse{
			Address: address.Hex(),
			High:    summary.High,
			Low:     summary.Low,
			Change:  summary.Change,
		})
	}

	c.JSON(200, response)

}
