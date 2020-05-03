package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Option func(s *server)

func WithPrometheus(path string) Option {
	return func(s *server) {
		s.prometheus = path
		if len(path) > 0 {
			s.router.GET(path, func(handler http.Handler) func(*gin.Context) {
				return func(c *gin.Context) {
					handler.ServeHTTP(c.Writer, c.Request)
				}
			}(promhttp.Handler()))
		}
	}
}

func WithFavicon(enable bool) Option {
	return func(s *server) {
		if enable {
			s.router.GET("/favicon.ico", func(c *gin.Context) {
				c.SecureJSON(http.StatusOK, gin.H{})
			})
		}
	}
}
