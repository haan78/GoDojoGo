package libecho

import (
	"strings"

	"github.com/labstack/echo/v5"
)

func RealIP(c *echo.Context) string {
	// Cloudflare provides CF-Connecting-IP
	if ip := c.Request().Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}
	// Fallbacks (if you use other proxies)
	if ip := c.Request().Header.Get("X-Forwarded-For"); ip != "" {
		// XFF can contain multiple IPs: client, proxy1, proxy2...
		// take the first
		if comma := strings.Index(ip, ","); comma > 0 {
			return strings.TrimSpace(ip[:comma])
		}
		return strings.TrimSpace(ip)
	}
	return c.RealIP()
}
