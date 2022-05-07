package services

import "github.com/gin-gonic/gin"

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) AlertSubscribeHandler(c *gin.Context) {

	var req SubscribeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, &Logger{
			Message: "Wrong body format",
		})
		return
	}

}
