// Package app configures and runs application.
package app

import (
	"context"
	"log"
	"net/http"
	"time"

	config "backend-service/config/marketplace"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type params struct {
	fx.In
	Cfg *config.Config
	Lc  fx.Lifecycle
}

// NewServer starts the http server.
func NewServer(params params) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())
	httpServer := &http.Server{
		Addr:              params.Cfg.HTTP.Host + ":" + params.Cfg.HTTP.Port,
		Handler:           engine,
		ReadHeaderTimeout: 2 * time.Second,
	}

	params.Lc.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				log.Println("Starting HTTPS server at " + httpServer.Addr)
				go func() {
					err := httpServer.ListenAndServe()
					if err != nil {
						panic(err)
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				log.Println("Stopping HTTPS server.")

				return httpServer.Shutdown(ctx)
			},
		},
	)

	return engine
}
