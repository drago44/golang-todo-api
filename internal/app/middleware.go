package app

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CORS middleware handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// Logger middleware logs HTTP requests
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

// Recovery middleware recovers from panics
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v", recovered)
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
	})
}

// RateLimit middleware for future rate limiting implementation
func RateLimit() gin.HandlerFunc {
	const (
		maxRequests = 100
		window      = time.Minute
	)

	type clientWindow struct {
		count      int
		windowEnds time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*clientWindow)
	)

	return gin.HandlerFunc(func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		cw, ok := clients[ip]
		if !ok || now.After(cw.windowEnds) {
			cw = &clientWindow{count: 0, windowEnds: now.Add(window)}
			clients[ip] = cw
		}
		if cw.count >= maxRequests {
			mu.Unlock()
			c.AbortWithStatusJSON(429, gin.H{
				"error":       "Too Many Requests",
				"retry_after": int(time.Until(cw.windowEnds).Seconds()),
			})
			return
		}
		cw.count++
		mu.Unlock()

		c.Next()
	})
}
