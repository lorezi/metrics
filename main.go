package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// computeSum generates a SHA256 hash of the body
func computeSum(body []byte) []byte {
	h := sha256.New()
	h.Write(body)
	hashed := hex.EncodeToString(h.Sum(nil))
	return []byte(hashed)
}

func main() {
	port := 8082
	readTimeout := time.Millisecond * 500
	writeTimeout := time.Millisecond * 500

	mux := http.NewServeMux()

	httpMetrics := NewHttpMetrics()
	instrumentedHash := InstrumentHandler(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.Body != nil {
			body, _ := ioutil.ReadAll(r.Body)
			fmt.Fprintf(w, "%s", computeSum(body))
		}
	}, httpMetrics)

	mux.HandleFunc("/hash", instrumentedHash)

	mux.Handle("/metrics", promhttp.Handler())

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: 1 << 20,
		Handler:        mux,
	}

	log.Printf("Listening on: %d", port)
	s.ListenAndServe()
}
