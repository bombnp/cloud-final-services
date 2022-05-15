package pair

import "github.com/gin-gonic/gin"

type Handler struct {
	s *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		s: s,
	}
}

func (h *Handler) GetPairs(c *gin.Context) {
	resp, err := h.s.GetPairs()

	if err != nil {
		c.JSON(400, &Logger{
			Message: err.Error(),
		})
		return
	}

	c.JSON(200, resp)
}
