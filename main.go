package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	availablePath = flag.String("available-path", "/readyz", "Allowed path for requests")
	addr          = flag.String("serve-addr", ":8080", "Addr to serve requests")
)

func main() {
	flag.Parse()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(strings.TrimLeft(r.URL.Path, "/"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if u.Path != *availablePath {
			http.Error(w, "404 not found", http.StatusNotFound)
			return
		}

		res, err := http.DefaultClient.Get(u.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer res.Body.Close()
		w.WriteHeader(res.StatusCode)
		_, _ = io.Copy(w, res.Body)
	})

	srv := &http.Server{
		Addr:    *addr,
		Handler: handler,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("signal ", Wait([]os.Signal{syscall.SIGTERM, syscall.SIGINT}).String())

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatal(err.Error())
	}
}

func Wait(signals []os.Signal) os.Signal {
	sig := make(chan os.Signal, len(signals))
	signal.Notify(sig, signals...)
	return <-sig
}
