package main

import (
	"bank-service/server"
	"bank-service/stg"
	"bank-service/stg_ex"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	local_storage := stg_ex.Newstorage()
	service := stg.NewBalanceService(local_storage)
	balanceServer := server.NewBalanceServer(*service)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: balanceServer.Mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Starting server on http://localhost:8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting the server: %v\n", err)
		}
	}()

	<-stop
	fmt.Println("\nReceived shutdown signal...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}
