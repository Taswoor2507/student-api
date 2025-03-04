package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Taswoor2507/student-api/internal/config"
)

func main() {
	// load config
	cfg := config.MustLoad()
	//database setup

	//setup router
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to student api"))
	})
	//  setup server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}
	slog.Info("Server startted ", slog.String("address", cfg.Addr))
	fmt.Printf("server started! on %s\n", cfg.HTTPServer.Addr)
	// server gracefully shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("failed to run server: %v", err)
		}
	}()

	<-done
	slog.Info("shutting down the server ....")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Failed to shutdow server ", slog.String("error", err.Error()))

	}
	// if err != nil {
	// 	slog.Error("Failed to shutdow server ", slog.String("error" , err.Error()))
	// }

	slog.Info("Server shutdown successfully")
}

// context , slog , signals
