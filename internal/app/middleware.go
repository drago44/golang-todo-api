package app

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSWithConfig returns a CORS middleware configured from application settings.
func CORSWithConfig(cfg *Config) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(cfg.Server.AllowedOrigins))
	for _, o := range cfg.Server.AllowedOrigins {
		allowed[o] = struct{}{}
	}
	allowCredentials := cfg.Server.AllowCredentials
	return gin.HandlerFunc(func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			if _, ok := allowed[origin]; ok {
				c.Header("Vary", "Origin")
				c.Header("Access-Control-Allow-Origin", origin)
			}
		} else {
			// No Origin header — do not set wildcard when credentials allowed
			if !allowCredentials {
				c.Header("Access-Control-Allow-Origin", "*")
			}
		}
		c.Header("Access-Control-Allow-Credentials", fmt.Sprintf("%t", allowCredentials))
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		c.Header("Access-Control-Max-Age", "600")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})
}

// Logger returns a middleware that logs HTTP requests in a custom format.
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\n",
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

// Recovery returns a middleware that recovers from panics and returns 500.
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		log.Printf("Panic recovered: %v", recovered)
		c.JSON(500, gin.H{
			"error": "Internal server error",
		})
	})
}

// RateLimit returns an in-memory rate limiter middleware with periodic cleanup.
func RateLimit() gin.HandlerFunc {
	const (
		maxRequests = 100
		window      = time.Minute
		evictAfter  = 10 * time.Minute
	)
	type clientWindow struct {
		count      int
		windowEnds time.Time
		lastSeen   time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*clientWindow)
	)
	// Cleanup goroutine
	go func() {
		Ticker := time.NewTicker(5 * time.Minute)
		defer Ticker.Stop()
		for range Ticker.C {
			mu.Lock()
			cut := time.Now().Add(-evictAfter)
			for ip, cw := range clients {
				if cw.lastSeen.Before(cut) {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	return gin.HandlerFunc(func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()

		mu.Lock()
		cw, ok := clients[ip]
		if !ok || now.After(cw.windowEnds) {
			cw = &clientWindow{count: 0, windowEnds: now.Add(window), lastSeen: now}
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
		cw.lastSeen = now
		mu.Unlock()

		c.Next()
	})
}
