package main

import (
	/*"log"

	"context"
	"fmt"
	"github.com/maffka123/gophermarktBonus/internal/app"
	"github.com/maffka123/gophermarktBonus/internal/config"
	"github.com/maffka123/gophermarktBonus/internal/handlers"
	"github.com/maffka123/gophermarktBonus/internal/storage"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"*/
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Print("starting...")
	/*	cfg, err := config.InitConfig()
		if err != nil {
			log.Fatalf("can't load config: %v", err)
		}

		logger, err := config.InitLogger(cfg.Debug, cfg.AppName)
		if err != nil {
			log.Fatalf("can't initialize zap logger: %v", err)
		}

		logger.Info("initializing the service...")

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		db, err := storage.InitDB(ctx, cfg, logger)
		if err != nil {
			logger.Fatal("Error initializing db", zap.Error(err))
		}

		r := handlers.BonusRouter(ctx, db, cfg.Key, logger)

		srv := &http.Server{Addr: cfg.Endpoint, Handler: r}

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		//See example here: https://pkg.go.dev/net/http#example-Server.Shutdown
		go func() {
			sig := <-quit
			logger.Info(fmt.Sprintf("caught sig: %+v", sig))
			if err := srv.Shutdown(ctx); err != nil {
				// Error from closing listeners, or context timeout:
				logger.Error("HTTP server Shutdown:", zap.Error(err))
			}
		}()

		statusTicker := time.NewTicker(time.Duration(30) * time.Second)
		go app.UpdateStatus(ctx, statusTicker.C, logger, db, cfg)

		logger.Info("Start serving on", zap.String("endpoint name", cfg.Endpoint))
		log.Fatal(srv.ListenAndServe())*/

	http.HandleFunc("/app/register", func(w http.ResponseWriter, req *http.Request) { fmt.Fprintf(w, "hello\n") })
	log.Fatal(http.ListenAndServe(":8080", nil))
	fmt.Print("done...")

}
