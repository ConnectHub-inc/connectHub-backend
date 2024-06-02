package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"

	"github.com/tusmasoma/connectHub-backend/config"
)

func main() {
	// .envファイルから環境変数を読み込む
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	var addr string
	flag.StringVar(&addr, "addr", ":8083", "tcp host:port to connect")
	flag.Parse()

	mainCtx, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()

	container, err := BuildContainer(mainCtx)
	if err != nil {
		log.Printf("Failed to build container: %v", err)
		return
	}

	/* ===== サーバの設定 ===== */
	err = container.Invoke(func(router *chi.Mux, config *config.ServerConfig) {
		srv := &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
			IdleTimeout:  config.IdleTimeout,
		}
		/* ===== サーバの起動 ===== */
		log.SetFlags(0)
		log.Println("Server running...")

		signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
		defer stop()

		go func() {
			if err = srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				log.Printf("Server failed: %v", err)
				return
			}
		}()

		<-signalCtx.Done()
		log.Println("Server stopping...")

		tctx, cancelShutdown := context.WithTimeout(context.Background(), config.GracefulShutdownTimeout)
		defer cancelShutdown()

		if err = srv.Shutdown(tctx); err != nil {
			log.Println("failed to shutdown http server", err)
		}
		log.Println("Server exited")
	})
	if err != nil {
		log.Printf("Failed to start server: %v", err)
		return
	}
}
