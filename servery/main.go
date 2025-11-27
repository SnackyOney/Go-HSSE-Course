package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const version string = "1.0.0"

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	response := map[string]string{
		"version": version,
	}

	json.NewEncoder(w).Encode(response)
}

func DecodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var request map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	encodedString, _ := request["inputString"].(string)

	decoded, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		fmt.Println("Error decoding:", err)
		return
	}

	response := map[string]string{
		"outputString": string(decoded),
	}

	json.NewEncoder(w).Encode(response)
}

func HardOpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	sleepTime := 10 + rand.Intn(11)
	time.Sleep(time.Duration(sleepTime) * time.Second)

	status := 500 + rand.Intn(100)

	response := map[string]int{
		"status": status,
	}

	json.NewEncoder(w).Encode(response)
}

func LaunchServer() {
	mux := http.NewServeMux()

	mux.HandleFunc("/version", VersionHandler)
	mux.HandleFunc("/decode", DecodeHandler)
	mux.HandleFunc("/hard-op", HardOpHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Starting server on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting the server: %v\n", err)
		}
	}()

	<-stop
	fmt.Println("\nReceived shutdown signal...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}
}

func main() {
	LaunchServer()
}
