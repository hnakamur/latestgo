package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hnakamur/latestgo"
)

func main() {
	queryTimeout := flag.Duration("query-timeout", 5*time.Second, "version query timeout")
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), *queryTimeout)
		defer cancel()

		w.Header().Set("Content-Type", "text/plain")

		ver, err := latestgo.Version(ctx)
		if err != nil {
			log.Printf("failed to get latest go version; %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		_, err = w.Write([]byte(ver))
		if err != nil {
			log.Printf("failed to write latest go version; %s", err)
			return
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}
	listenAddr := ":" + port
	srv := http.Server{
		Addr: listenAddr,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigC := make(chan os.Signal, 1)
		signal.Notify(sigC, os.Interrupt, syscall.SIGTERM)
		<-sigC

		// We received an interrupt or SIGTERM signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
