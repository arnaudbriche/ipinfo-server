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
	mux.HandleFunc("/ipinfo", ipInfoHandler)
	mux.HandleFunc("/lookup", lookupHandler)
	mux.HandleFunc("/lookupsrv", lookupSRVHandler)
	mux.HandleFunc("/dial", dialHandler)

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

func lookupHandler(w http.ResponseWriter, r *http.Request) {
	var hostname = r.URL.Query().Get("hostname")

	if len(hostname) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ips, err := net.LookupHost(hostname)

	if err != nil && len(ips) == 0 {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := json.NewEncoder(w).Encode(ips); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type srvTarget struct {
	Hostname string
	Port     int16
	IPS      []string
	Error    string
}

type lookupSRVResponse struct {
	CNAME   string
	Targets []srvTarget
}

func lookupSRVHandler(w http.ResponseWriter, r *http.Request) {
	var (
		domain  = r.URL.Query().Get("domain")
		service = r.URL.Query().Get("service")
	)

	if len(domain) == 0 || len(service) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cname, addrs, err := net.LookupSRV(service, "tcp", domain)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var resp lookupSRVResponse
	resp.CNAME = cname

	for _, addr := range addrs {
		var target srvTarget

		target.Hostname = addr.Target
		target.Port = int16(addr.Port)

		ips, err := net.LookupHost(addr.Target)

		if err != nil {
			target.Error = err.Error()
			resp.Targets = append(resp.Targets, target)
			continue
		}

		target.IPS = ips
		resp.Targets = append(resp.Targets, target)
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func dialHandler(w http.ResponseWriter, r *http.Request) {
	var address = r.URL.Query().Get("address")

	conn, err := net.Dial("tcp", address)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer conn.Close()

	var resp = make(map[string]interface{})

	resp["LocalAddr"] = conn.LocalAddr()
	resp["RemoteAddr"] = conn.RemoteAddr()

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
