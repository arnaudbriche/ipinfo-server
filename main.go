package main

import (
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"os"
)

var logger = slog.New(slog.NewTextHandler(os.Stderr, nil))

func main() {
	var bindAddr = ":8080"

	if len(os.Getenv("BIND_ADDR")) > 0 {
		bindAddr = os.Getenv("BIND_ADDR")
	}

	var mux = http.NewServeMux()
	mux.HandleFunc("/", ipInfoHandler)
	mux.HandleFunc("/resolve", resolveHandler)

	var server = &http.Server{
		Addr:    bindAddr,
		Handler: mux,
	}

	logger.Info("HTTP server running", "addr", bindAddr)

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func ipInfoHandler(w http.ResponseWriter, r *http.Request) {
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

func resolveHandler(w http.ResponseWriter, r *http.Request) {
	var hostname = r.URL.Query().Get("hostname")

	if len(hostname) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ips, err := net.LookupIP(hostname)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(ips); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
