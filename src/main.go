package main

import (
	"demo/src/config"
	"demo/src/engine"
	"demo/src/logger"
	"demo/src/memcache"
	"demo/src/pkg/httpx"
	"demo/src/router"
	"demo/src/storage"
	"fmt"
	"github.com/jinzhu/configor"
	"os"
	"os/signal"
	"syscall"
)

var c *config.Config

func main() {
	code := 1
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	err := configor.Load(&c, "etc/config.yml")

	// init logger
	if err = logger.InitLogger(c.Log); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	// init database
	if err = storage.InitDB(c.Db); err != nil {
		logger.Error("database init failed", err)
	}

	// init redis
	_, err = storage.InitRedis(c.Redis)
	if err != nil {
		logger.Error("redis init failed", err)
	}

	memcache.SyncTradingPairs()
	go engine.Start()

	// init http
	httpEngine := router.New(c.Server)
	httpx.Init(c.Server, httpEngine)

	logger.Info("started")

EXIT:
	for {
		sig := <-sc
		fmt.Println("received signal:", sig.String())
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			code = 0
			break EXIT
		default:
			break EXIT
		}
	}

	fmt.Println("server exited")
	os.Exit(code)

}
