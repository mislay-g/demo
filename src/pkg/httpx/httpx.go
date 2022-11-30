package httpx

import (
	"demo/src/config"
	"demo/src/logger"
	"fmt"
	"net/http"
	"time"
)

func Init(cfg *config.ServerConfig, handler http.Handler) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeout) * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("http server listening on: %s", addr))

		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Error("http server started failed", err)
			panic(err)
		}
	}()

	//return func() {
	//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(cfg.ShutdownTimeout))
	//	defer cancel()
	//
	//	srv.SetKeepAlivesEnabled(false)
	//	if err := srv.Shutdown(ctx); err != nil {
	//		fmt.Println("cannot shutdown http server:", err)
	//	}
	//
	//	select {
	//	case <-ctx.Done():
	//		fmt.Println("http exiting")
	//	default:
	//		fmt.Println("http server stopped")
	//	}
	//}
}
