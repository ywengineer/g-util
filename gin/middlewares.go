package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/ywengineer/g-util/api"
	"github.com/ywengineer/g-util/cache"
	"net/http"
	"strconv"
	"time"
)

func MiddlewareZipKin(tracer *zipkin.Tracer) gin.HandlerFunc {
	return func(context *gin.Context) {
		span, ctx := tracer.StartSpanFromContext(context.Request.Context(), context.FullPath(), zipkin.Kind(model.Server))
		defer span.Finish()
		context.Request = context.Request.WithContext(ctx)

		zipkin.TagHTTPMethod.Set(span, context.Request.Method)
		zipkin.TagHTTPPath.Set(span, context.Request.URL.Path)
		zipkin.TagHTTPRequestSize.Set(span, fmt.Sprintf("%d", context.Request.ContentLength))
		//
		context.Next()
		//
		if statusCode := context.Writer.Status(); statusCode < 200 || statusCode > 299 {
			zipkin.TagHTTPStatusCode.Set(span, strconv.Itoa(context.Writer.Status()))
			zipkin.TagHTTPResponseSize.Set(span, strconv.Itoa(context.Writer.Size()))
		}
	}
}

func MiddlewareToken(transport *zipkinhttp.Client,
	enableCache bool,
	tokenName, tokenUrl string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get(tokenName)
		if len(token) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, api.Result{
				Code:    api.ERR400,
				Message: "Missing Token",
			})
			return
		}
		// find cache
		var ci interface{}
		ok := false
		if enableCache {
			ci, ok = cache.Get(token)
		}
		// not found
		if !ok {
			// "https://console.linkfungame.com/lfact/api/user"
			req, err := http.NewRequestWithContext(c.Request.Context(), "GET", tokenUrl, nil)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.Result{
					Code:    api.ERR500,
					Message: "StatusInternalServerError" + err.Error(),
				})
				return
			} else {
				q := req.URL.Query()
				q.Add("token", token)
				req.URL.RawQuery = q.Encode()
			}
			res, err := transport.DoWithAppSpan(req, req.URL.Path)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, api.Result{
					Code:    api.ERR500,
					Message: "StatusInternalServerError: " + err.Error(),
				})
				return
			}
			////////
			result := &api.Result{}
			//
			if err := result.ReadFrom(res); err != nil || !result.IsSuccess() || result.Data == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, api.Result{
					Code:    api.ERR500,
					Message: "invalid token",
				})
				if enableCache {
					cache.Set(token, result, 25*time.Minute)
				}
				return
			}
			// cache
			if enableCache {
				cache.Set(token, result, 1*time.Hour)
			}
		} else {
			t, ok := ci.(*api.Result)
			// invalid token
			if !ok || !t.IsSuccess() || t.Data == nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, api.Result{
					Code:    api.ERR500,
					Message: "invalid token",
				})
				return
			}
		}
	}
}
