package app

import (
	"log"
	"net/http"
	"time"

	"github.com/drago44/golang-todo-api/internal/router"
	"github.com/drago44/golang-todo-api/internal/todos"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

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
	if err := container.Provide(func() *gin.Engine {
		gin.SetMode(gin.ReleaseMode)
		engine := gin.New()
		engine.Use(Logger(), Recovery(), CORS(), RateLimit())
		// TODO: Add trusted proxies
		if err := engine.SetTrustedProxies([]string{}); err != nil {
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

	if err := container.Provide(router.New); err != nil {
		log.Fatal(err)
	}

	if err := container.Invoke(func(router *router.Router, cfg *Config) {
		addr := cfg.Server.Host + ":" + cfg.Server.Port
		log.Printf("Server starting on %s", addr)
		srv := &http.Server{
			Addr:              addr,
			Handler:           router.GetEngine(),
			ReadTimeout:       5 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      10 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    1 << 20,
		}
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}); err != nil {
		log.Fatal(err)
	}

}
