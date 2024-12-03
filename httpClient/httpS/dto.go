/**
 * Package httpS
 * @Author fengfeng.mei <fengfeng.mei@baishan.com>
 * @Date 2024/12/3 14:30
 */

package httpS

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/pprof"
	"sync"
	"time"
)

type (
	HttpServers struct {
		server *http.Server
		perf   sync.Map
	}
)

var (
	instance *HttpServers
	once     sync.Once
)

func GetInstance() *HttpServers {
	once.Do(func() {
		instance = &HttpServers{}
		gin.SetMode("debug")
		r := gin.New()
		// 注册restful路由
		RegisterRestfulRoute(r, instance)
		instance.server = &http.Server{
			Addr:           "127.0.0.1:8989",
			Handler:        r,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
	})
	return instance
}

func RegisterRestfulRoute(r *gin.Engine, s *HttpServers) {
	pprofR := r.Group("/debug/pprof")
	pprofR.GET("/", gin.WrapF(pprof.Index))
	pprofR.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	pprofR.GET("/profile", gin.WrapF(pprof.Profile))
	pprofR.POST("/symbol", gin.WrapF(pprof.Symbol))
	pprofR.GET("/symbol", gin.WrapF(pprof.Symbol))
	pprofR.GET("/trace", gin.WrapF(pprof.Trace))
	pprofR.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
	pprofR.GET("/block", gin.WrapH(pprof.Handler("block")))
	pprofR.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
	pprofR.GET("/heap", gin.WrapH(pprof.Handler("heap")))
	pprofR.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
	pprofR.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
}

func (hs *HttpServers) Start() error {
	go func() {
		err := hs.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("start-fail", "addr", hs.server.Addr, "err", err)
		}
	}()
	return nil
}
