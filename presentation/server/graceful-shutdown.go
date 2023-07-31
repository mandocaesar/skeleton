package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/machtwatch/catalyst-go-skeleton/infrastructure/config"
	"github.com/machtwatch/catalystdk/go/log"
	"github.com/machtwatch/catalystdk/go/trace"
)

// startServerWithGracefulShutdown start the application server with graceful shutdown.\
//
// On graceful shutdown, it shuts down the server without interrupting any
// active connections.
func startServerWithGracefulShutdown(r *chi.Mux, tracer *trace.TracerSet) {
	addr := fmt.Sprintf("0.0.0.0:%d", config.APP_PORT)
	server := &http.Server{Addr: addr, Handler: r}

	// Create server context
	serverCtx, cancelServerCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		graceful_timeout := config.SERVER_GRACEFUL_SHUTDOWN_TIMEOUT_S

		shutDownCtx, cancel := context.WithTimeout(serverCtx, time.Duration(graceful_timeout)*time.Second)
		tpCtx, tpCancel := context.WithTimeout(serverCtx, time.Duration(graceful_timeout)*time.Second)

		defer cancel()
		defer tpCancel()

		go func() {
			<-shutDownCtx.Done()
			if shutDownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out, forcing exit..")
			}
		}()

		log.Info("shutting down tracer with timeout: %v seconds", graceful_timeout)

		err := tracer.Tracer.Shutdown(tpCtx)
		if err != nil {
			log.Fatalf("error on shutting down tracer gracefully: %v", err)
		}

		err = tracer.Metric.Shutdown(tpCtx)
		if err != nil {
			log.Fatalf("error on shutting down middleware tracer gracefully: %v", err)
		}

		log.Info("shutting down server with timeout: %v seconds", graceful_timeout)

		err = server.Shutdown(shutDownCtx)
		if err != nil {
			log.Fatalf("error on shutting down gracefully: %v", err)
		}

		cancelServerCtx()
	}()

	log.Infof("starting %v in port: %v", config.APP_NAME, config.APP_PORT)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("error on starting up server: %v", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	log.Info("server is shut down!")
}
