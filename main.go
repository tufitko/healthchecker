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
	"regexp"
	"strings"
	"syscall"
)

var (
	availablePath = flag.String("available-path", "/readyz", "Allowed path for requests")
	addr          = flag.String("serve-addr", ":8080", "Addr to serve requests")

	allowedHosts = flag.String("allowed-hosts", "", "Allowed hosts for requests (regex)")
)

func main() {
	flag.Parse()

	isAllowed := allowedHostsNoop
	if *allowedHosts != "" {
		var err error
		isAllowed, err = allowedHostsRegex(*allowedHosts)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(strings.TrimLeft(r.URL.Path, "/"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if u.Path != *availablePath || !isAllowed(u.Hostname()) {
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

func allowedHostsNoop(_ string) bool {
	return true
}

func allowedHostsRegex(allowedRegex string) (func(string) bool, error) {
	r, err := regexp.Compile(allowedRegex)
	if err != nil {
		return nil, err
	}

	return func(s string) bool {
		return r.MatchString(s)
	}, nil
}
