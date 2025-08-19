package app

import (
	"log"

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
	container.Provide(func() *Config { return cfg })
	container.Provide(func() *gorm.DB { return db })
	container.Provide(func() *gin.Engine {
		engine := gin.New()
		engine.Use(Logger(), Recovery(), CORS(), RateLimit())
		// TODO: Add trusted proxies
		if err := engine.SetTrustedProxies([]string{}); err != nil {
			log.Fatal(err)
		}
		return engine
	})

	for _, module := range []func(*dig.Container) error{
		todos.Module,
	} {
		if err := module(container); err != nil {
			log.Fatal(err)
		}
	}

	container.Provide(router.New)

	container.Invoke(func(router *router.Router, cfg *Config) {
		log.Printf("Server started on %s:%s", cfg.Server.Host, cfg.Server.Port)
		err := router.GetEngine().Run(cfg.Server.Host + ":" + cfg.Server.Port)
		if err != nil {
			log.Fatal(err)
		}
	})

}
