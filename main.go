package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	cert := os.Getenv("TLS_CERT")
	key := os.Getenv("TLS_KEY")

	errLog := log.New(os.Stdout, "ERR: ", log.LstdFlags)

	hostname, err := os.Hostname()
	if err != nil {
		errLog.Printf("Failed to get hostname: %v\n", err)
	}

	var data Response

	go func() {
		client := &http.Client{
			Transport: &http.Transport{},
		}

		resp, err := client.Get("https://ifconfig.me/all.json")
		if err != nil {
			errLog.Printf("Failed to get response: %v\n", err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			errLog.Printf("Failed to read response body: %v\n", err)
		}

		if json.Unmarshal(body, &data) != nil {
			errLog.Printf("Failed to unmarshal response body: %v\n", err)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Hello World from %v to %v via %v. Egress %v via %v", hostname, r.RemoteAddr, r.Proto, data.egress(), data.Fowarded)))
	})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	addr := "0.0.0.0:8080"

	if cert != "" && key != "" {
		log.Fatal(http.ListenAndServeTLS(addr, cert, key, nil))
	} else {
		log.Fatal(http.ListenAndServe(addr, nil))
	}
}
