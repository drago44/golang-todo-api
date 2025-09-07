package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	docs "github.com/drago44/golang-todo-api/docs/swagger"
	"github.com/drago44/golang-todo-api/internal/router"
	"github.com/drago44/golang-todo-api/internal/todos"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

// Run performs application bootstrapping and starts the HTTP server.
func Run() {
	cfg, err := Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := Init(&cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	if err := Migrate(db); err != nil {
		log.Fatal(err)
	}

	container := dig.New()

	// Provide core singletons
	if err := container.Provide(func() *Config { return cfg }); err != nil {
		log.Fatal(err)
	}
	if err := container.Provide(func() *gorm.DB { return db }); err != nil {
		log.Fatal(err)
	}

	if err := container.Provide(func(cfg *Config) *gin.Engine {
		// Mode
		mode := cfg.Server.GinMode
		if mode == "" {
			mode = gin.ReleaseMode
		}
		gin.SetMode(mode)

		engine := gin.New()
		if cfg.Server.EnableLogger {
			engine.Use(Logger())
		}
		engine.Use(Recovery(), CORSWithConfig(cfg))
		if cfg.Server.EnableRateLimit {
			engine.Use(RateLimit())
		}
		// Trusted proxies
		if err := engine.SetTrustedProxies(cfg.Server.TrustedProxies); err != nil {
			log.Fatal(err)
		}
		return engine
	}); err != nil {
		log.Fatal(err)
	}

	for _, module := range []func(*dig.Container) error{
		todos.Module,
	} {
		if err := module(container); err != nil {
			log.Fatal(err)
		}
	}

	if err := container.Provide(func(engine *gin.Engine, todoHandler *todos.TodoHandler, cfg *Config) *router.Router {
		return router.New(engine, todoHandler, cfg.Server.EnableSwagger)
	}); err != nil {
		log.Fatal(err)
	}

	if err := container.Invoke(func(router *router.Router, cfg *Config) {
		addr := cfg.Server.Host + ":" + cfg.Server.Port

		// Determine the public scheme from config; fallback by port if not set
		protocol := cfg.Server.PublicScheme
		if protocol == "" {
			if cfg.Server.Port == "443" {
				protocol = "https"
			} else {
				protocol = "http"
			}
		}

		// Configure swagger metadata at runtime
		docs.SwaggerInfo.BasePath = "/api/v1"
		docs.SwaggerInfo.Host = addr
		docs.SwaggerInfo.Schemes = []string{protocol}

		// Form the URL for clickability
		url := protocol + "://" + addr

		log.Printf("🚀 Server starting on %s", url)
		if cfg.Server.EnableSwagger {
			log.Printf("📖 API Documentation: %s/swagger/index.html", url)
		}
		log.Printf("💚 Health Check: %s/health", url)

		// Start the server
		srv := &http.Server{
			Addr:              addr,
			Handler:           router.GetEngine(),
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    1 << 20,
		}

		// Start the server in a goroutine
		go func() {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}()

		// Shutdown the server gracefully
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		<-quit
		log.Printf("📦 Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server forced to shutdown: %v", err)

		}
	}); err != nil {
		log.Fatal(err)
	}
}
