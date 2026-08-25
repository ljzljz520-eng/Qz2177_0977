package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"coursechain/api"
	"coursechain/store"
	"coursechain/workflow"
)

func main() {
	path := flag.String("db", "coursechain.db", "bbolt database path")
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
	database, err := store.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	service := workflow.NewService(database)
	server := &http.Server{Addr: *addr, Handler: api.NewRouter(service), ReadHeaderTimeout: 3 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second, IdleTimeout: 30 * time.Second}
	go func() {
		if serveErr := server.ListenAndServe(); serveErr != nil && serveErr != http.ErrServerClosed {
			log.Printf("server stopped: %v", serveErr)
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "shutdown: %v\n", err)
	}
}
