package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"google.golang.org/appengine"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	mux := chi.NewMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/ip", ipAddressHandler)
	mux.HandleFunc("/headers", printHeadersHandler)

	mux.NotFound(notFoundHandler)
	mux.Mount("/debug", middleware.Profiler())
	mux.Handle("/metrics", promhttp.Handler())

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func ipAddressHandler(w http.ResponseWriter, r *http.Request) {
	var ip string
	if appengine.IsAppEngine() {
		ip = r.Header.Get("X-Appengine-User-Ip")
	} else {
		var err error
		ip, _, err = net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			w.WriteHeader(500)
			_, _ = fmt.Fprintf(w, "ip: %s is not IP:port", r.RemoteAddr)
			return
		}
	}

	_, _ = fmt.Fprintf(w, "%s\n", ip)
}

func printHeadersHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		if len(v) > 0 {
			_, _ = fmt.Fprintf(w, "%s: %s\n", k, v[0])
		}
	}
}

// notFoundHandler responds to Paths that are not otherwise supported
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprint(w, "Hello, World!")
}
