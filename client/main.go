package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func GetVersion(client http.Client) {

	resp, err := client.Get("http://localhost:8080/version")
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var versionInfo map[string]string
	err = json.NewDecoder(resp.Body).Decode(&versionInfo)
	if err != nil {
		fmt.Printf("JSON decode failed: %v\n", err)
		return
	}

	fmt.Printf("version: %s\n", versionInfo["version"])
}

func GetDecodedString(client http.Client, str string) {

	data := map[string]interface{}{
		"inputString": str,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("JSON marshal error: %v\n", err)
		return
	}

	resp, err := client.Post("http://localhost:8080/decode", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("JSON decode failed: %v\n", err)
		return
	}

	fmt.Printf("decoded string: %s\n", string(result["outputString"]))
}

func GetHardOp(client http.Client) {

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/hard-op", nil)
	if err != nil {
		fmt.Printf("Request creation failed: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var statusInfo map[string]int
	err = json.NewDecoder(resp.Body).Decode(&statusInfo)
	if err != nil {
		fmt.Printf("JSON decode failed: %v\n", err)
		return
	}

	fmt.Printf("status: %v\n", statusInfo["status"])
}

func main() {
	client := http.Client{}
	GetVersion(client)
	GetDecodedString(client, "SGVsbG8=")
	GetHardOp(client)
}
