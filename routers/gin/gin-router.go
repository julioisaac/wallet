package gin

import (
	"fmt"
	v3 "github.com/SkyAPM/go2sky-plugins/gin/v3"
	"github.com/gin-gonic/gin"
	"github.com/julioisaac/daxxer-api/metrics"
	"github.com/julioisaac/daxxer-api/routers"
	"net/http"
)

type ginRouter struct{}

var (
	ginDispatcher = gin.Default()
)

func NewGinRouter() routers.Router {
	gin.SetMode(gin.ReleaseMode)
	return &ginRouter{}
}

func (g *ginRouter) SetupTracer() {
	tracer := metrics.Metric().GetTracer()
	ginDispatcher.Use(v3.Middleware(ginDispatcher, tracer))
}

func (g *ginRouter) GET(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	ginDispatcher.GET(uri, gin.WrapF(f))
}

func (g *ginRouter) POST(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	ginDispatcher.POST(uri, gin.WrapF(f))
}

func (g *ginRouter) PUT(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	ginDispatcher.PUT(uri, gin.WrapF(f))
}

func (g *ginRouter) DELETE(uri string, f func(w http.ResponseWriter, r *http.Request)) {
	ginDispatcher.DELETE(uri, gin.WrapF(f))
}

func (g *ginRouter) SERVE(port string) {
	err := ginDispatcher.Run(port)
	if err != nil {
		fmt.Printf("Error trying running Gin")
		return
	}
}