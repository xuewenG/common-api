package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/samber/do/v2"
)

type IPHandler struct {
	*Handler
}

func NewIPHandler(i do.Injector) (*IPHandler, error) {
	h, err := newHandler(i)
	if err != nil {
		return nil, err
	}

	return &IPHandler{
		Handler: h,
	}, nil
}

func (h *IPHandler) Echo(c *gin.Context) {
	ip := h.getRealIP(c)

	format := c.Query("format")
	if format == "json" {
		c.JSON(http.StatusOK, gin.H{"ip": ip})
		return
	}

	c.String(http.StatusOK, ip)
}

func (h *IPHandler) getRealIP(c *gin.Context) string {
	xRealIP := c.GetHeader("x-real-ip")
	xForwardedFor := c.GetHeader("x-forwarded-for")
	clientIP := c.ClientIP()

	h.logger.Info().Msgf("x-real-ip: %s, x-forwarded-for: %s, clientIP: %s", xRealIP, xForwardedFor, clientIP)

	if xRealIP != "" {
		return xRealIP
	}

	if xForwardedFor != "" {
		ips := strings.Split(xForwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	return clientIP
}
