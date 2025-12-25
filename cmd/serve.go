package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-playground/validator"
	"github.com/spl3g/tcpinger/handlers"
)

func handleServe(host string, port int) error {
	sigCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	validate := validator.New()

	s := handlers.NewServer(validate)
	addr := net.JoinHostPort(host, fmt.Sprint(port))
	server := &http.Server{
		Addr:         addr,
		Handler:      s,
		BaseContext:  func(_ net.Listener) context.Context { return sigCtx },
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
		IdleTimeout:  time.Second,
	}

	errChan := make(chan error)

	go func() {
		defer cancel()
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("HTTP server error: %s", err)
		}
	}()

	log.Printf("Started HTTP server on %s\n", addr)
	select {
	case <-sigCtx.Done():
		log.Printf("Stopping server")
		ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			return err
		}
	case err := <-errChan:
		return err
	}

	return nil
}
