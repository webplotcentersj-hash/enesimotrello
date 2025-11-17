package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get allowed origin from config or environment
		origin := c.Request.Header.Get("Origin")
		allowedOrigin := os.Getenv("CORS_ORIGIN")
		
		// If CORS_ORIGIN is set, use it; otherwise allow all origins
		if allowedOrigin == "" || allowedOrigin == "*" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else {
			// Support multiple origins (comma-separated)
			allowedOrigins := strings.Split(allowedOrigin, ",")
			for _, allowed := range allowedOrigins {
				allowed = strings.TrimSpace(allowed)
				if origin == allowed {
					c.Header("Access-Control-Allow-Origin", origin)
					c.Header("Access-Control-Allow-Credentials", "true")
					break
				}
			}
			// If no match found but we have a single allowed origin, use it
			if c.Header("Access-Control-Allow-Origin") == "" && len(allowedOrigins) == 1 {
				c.Header("Access-Control-Allow-Origin", strings.TrimSpace(allowedOrigins[0]))
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}

		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Anonymous-User-Id")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Header("Access-Control-Expose-Headers", "Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
