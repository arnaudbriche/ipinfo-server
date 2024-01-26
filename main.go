package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func main() {
	var bindAddr = ":8080"

	if len(os.Getenv("BIND_ADDR")) > 0 {
		bindAddr = os.Getenv("BIND_ADDR")
	}

	var server = &http.Server{
		Addr:    bindAddr,
		Handler: http.HandlerFunc(handler),
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	var client http.Client

	req, err := http.NewRequest("GET", "https://ipinfo.io", nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := client.Do(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var apiResp map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(apiResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
