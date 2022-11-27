package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"pasto/httpserver"
	. "pasto/logger"
	"syscall"
	"time"
)

var listenPort string

func main() {
	flag.StringVar(&listenPort, "l", "8088", "Port to listen")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", httpserver.RootHandler)
	httpServ := &http.Server{
		Addr:              ":" + listenPort,
		Handler:           mux,
		ReadHeaderTimeout: 60 * time.Second,
	}

	go func() {
		Logger.Info().Err(httpServ.ListenAndServe()).
			Msg("Initiating shutdown")
	}()
	Logger.Info().Msgf("Listening on port %s", listenPort)

	// Graceful shutdown for HTTP server
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	Logger.Info().Msg("HTTP server stopping")
	defer cancel()
	Logger.Info().Err(httpServ.Shutdown(ctx)).Msg("HTTP server stopped")
}
