package subscribe

import "github.com/gin-gonic/gin"

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) PostAlertSubscribe(c *gin.Context) {

	var req AlertSubscribeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, &Logger{
			Message: "Wrong body format",
		})
		return
	}

	if err := h.s.AlertSubscribe(req.ServerId, req.PoolAddress, req.ChannelId); err != nil {
		c.JSON(400, &Logger{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, &Logger{
		Message: "OK",
	})

}
