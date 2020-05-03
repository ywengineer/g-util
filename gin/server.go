package gin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/ywengineer/g-util/api"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type server struct {
	engine     *http.Server
	router     *gin.Engine
	prometheus string
	logger     *zap.Logger
}

func NewGinServer(port int, logger *zap.Logger, opts ...Option) *server {
	r := gin.Default()
	s := &server{
		engine: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      r,
			ReadTimeout:  2 * time.Second,
			WriteTimeout: 3 * time.Second,
		},
		logger: logger,
	}
	s.router = r
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(s)
		}
	}
	return s
}

func (s *server) Router() *gin.Engine {
	return s.router
}

func (s *server) Run() error {
	//
	s.router.NoRoute(func(c *gin.Context) {
		c.SecureJSON(http.StatusBadRequest, api.Result{
			Code:    api.ERR400,
			Message: "illegal",
			Data:    nil,
		})
		return
	})
	err := s.engine.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("run gin service failed.", zap.Error(err))
	}
	return err
}

func (s *server) Terminate() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.engine.Shutdown(ctx); err != nil {
		s.logger.Error("timeout to shutdown gin service", zap.Error(err))
	}
	s.logger.Info("gin service is terminated.")
}
