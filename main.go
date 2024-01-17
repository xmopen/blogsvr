package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xmopen/golib/pkg/xgoroutine"

	"github.com/xmopen/golib/pkg/xlogging"

	"github.com/xmopen/blogsvr/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/xmopen/blogsvr/internal/endpoint"
)

type app struct {
	engine *gin.Engine
	apiSvr *http.Server
	cancel context.CancelFunc
	close  chan error
	xlog   *xlogging.Entry
}

// init 初始化svr.
func (a *app) init(ctx context.Context) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)
	go func() {
		select {
		case r := <-sigs:
			a.close <- fmt.Errorf("syscall:[%+v]\n", r)
		}
	}()

	endpoint.Init(a.engine)
	a.run(ctx)
}

// run 运行svr.
func (a *app) run(ctx context.Context) {
	if err := a.apiSvr.ListenAndServe(); err != nil {
		a.close <- err
	}
}

func (a *app) quit() {
	select {
	case err := <-a.close:
		fmt.Println("svr done because err:" + err.Error())
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	r := gin.New()
	addr := config.Config().GetString("server.blogsvr.http.addr")
	app := &app{
		engine: r,
		apiSvr: &http.Server{
			Addr:              addr,
			Handler:           r,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      5 * time.Second,
		},
		cancel: cancel,
		close:  make(chan error, 1), // 容量为1不阻塞.
		xlog:   xlogging.Tag("blogsvr.main"),
	}
	xgoroutine.SafeGoroutine(func() {
		authSvrHostName := os.Getenv("AUTHSVR_SERVICE_HOST")
		authSvrHostPort := os.Getenv("AUTHSVR_SERVICE_PORT_AUTHSVRRPCPORT")
		app.xlog.Infof("authsvr service inner kubernetes ip:[%+v] port:[%+v]", authSvrHostName, authSvrHostPort)
	})
	app.xlog.Infof("http server running in addr:[%+v]", addr)
	app.init(ctx)
	app.quit()
}
