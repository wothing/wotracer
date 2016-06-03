package main

import (
	"net"
	"net/http"
	"net/url"
	"time"

	"sourcegraph.com/sourcegraph/appdash"
	"sourcegraph.com/sourcegraph/appdash/traceapp"

	"github.com/wothing/log"
)

func main() {
	// Create memory store for the internal data storage, and evicting
	// data after 10 minutes
	memStore := appdash.NewMemoryStore()
	store := &appdash.RecentStore{
		MinEvictAge: 10 * time.Minute,
		DeleteStore: memStore,
	}

	// Start the Appdash web UI on port 8700
	url, err := url.Parse("http://localhost:8700")
	if err != nil {
		log.Fatalf("failed to parsing url, error: %v", err)
	}
	tapp, err := traceapp.New(nil, url)
	if err != nil {
		log.Fatalf("failed to create trace app, error: %v", err)
	}
	tapp.Store = store
	tapp.Queryer = memStore

	log.Infof("AppDash web UI running on HTTP :8700")
	go func() {
		log.Fatal(http.ListenAndServe(":8700", tapp))
	}()

	// Start the Appdash trace collector server on port 1726
	l, err := net.Listen("tcp", ":1726")
	if err != nil {
		log.Fatalf("fail to listen to port, error: %v", err)
	}

	log.Infof("AppDash collector listening on :1726")
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	cs.Debug = false
	cs.Trace = false
	cs.Start()
}
