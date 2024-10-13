package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"statefulset/base"
	"statefulset/cmds/server/context"
	"statefulset/cmds/server/utils"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type NewRequestID func() (int64, string)
type HandlerFunc func(ctx *context.Context) error

type ApiServer struct {
	router       *gin.Engine
	Logger       *logrus.Logger
	NewRequestID NewRequestID
	workers      int64
	Listen       *base.ServerSetting
}

func NewApiServer(logger *logrus.Logger) *ApiServer {
	router := gin.New()
	router.Use(gin.Recovery())
	return &ApiServer{
		router:       router,
		Logger:       logger,
		NewRequestID: utils.NewRequestID,
		Listen:       &base.Context.Setting.Server,
	}
}

func (api *ApiServer) Start() error {
	server := &http.Server{
		Addr:           api.Listen.Addr,
		Handler:        api.router,
		ReadTimeout:    api.Listen.ReadTimeout,
		WriteTimeout:   api.Listen.WriteTimeout,
		IdleTimeout:    api.Listen.IdleTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Println(api.Listen.Addr, api.Listen.ReadTimeout, time.Second*time.Duration(api.Listen.ReadTimeout))

	return server.ListenAndServe()
}
func (api *ApiServer) increaseWorkers() {
	atomic.AddInt64(&api.workers, 1)
}

func (api *ApiServer) decreaseWorkers() {
	atomic.AddInt64(&api.workers, ^0)
}

func (api *ApiServer) AddRoute(method, relativePath string, handlers HandlerFunc) gin.IRoutes {
	handler := func(ctx *gin.Context) {
		var err error

		start := time.Now()
		api.increaseWorkers()

		RequestIDInt64, requestIDStr := api.NewRequestID()

		reqURL := ctx.Request.Method + " " + ctx.Request.Host + ctx.Request.URL.String()

		api.Logger.Infof("Request Received: %s", reqURL)
		api.Logger.Infof("Request headers: %v", ctx.Request.Header)

		c := context.NewContext(ctx)
		c.Logger = api.Logger
		c.RequestIDInt64 = RequestIDInt64
		c.RequestIDStr = requestIDStr

		defer func() {
			if r := recover(); r != nil {
				api.Logger.Errorf("Receive panic: %v", r)
				api.Logger.Errorf("Stack: %s", string(debug.Stack()))
				err = fmt.Errorf("Panic route")
			}
			if err != nil {
				api.Logger.Errorf("recover panic err: %v", err)
			}
			// Discard useless body
			c.Discard()

			latency := time.Now().Sub(start).Seconds() * 1000
			if c.StatusCode >= 500 {
				api.Logger.Errorf("Request Completed [%d]: %s, exec_time: [%.5fms]", c.StatusCode, reqURL, latency)
			} else {
				api.Logger.Infof("Request Completed [%d]: %s, exec_time: [%.5fms]", c.StatusCode, reqURL, latency)
			}
			// Release Request Context
			c.Free()
			api.decreaseWorkers()
			return
		}()
		err = handlers(c)
	}
	return api.router.Handle(method, relativePath, handler)

}
func (api *ApiServer) AddRouteWithoutM(method, relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes {
	return api.router.Handle(method, relativePath, handlers...)
}
