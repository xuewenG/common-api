package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type HealthHandler struct{}

func NewHealthHandler(_ do.Injector) (*HealthHandler, error) {
	return &HealthHandler{}, nil
}

func (h *HealthHandler) Check(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
