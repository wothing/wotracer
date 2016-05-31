package main

import (
	"net"
	"net/http"
	"net/url"

	"sourcegraph.com/sourcegraph/appdash"
	"sourcegraph.com/sourcegraph/appdash/traceapp"

	"github.com/wothing/log"
)

func main() {
	// Create a default InfluxDB configuration
	conf, err := appdash.NewInfluxDBConfig()
	if err != nil {
		log.Fatalf("failed to create influxdb config, error: %v", err)
	}

	// Enable InfluxDB server HTTP basic auth
	conf.Server.HTTPD.AuthEnabled = true
	conf.AdminUser = appdash.InfluxDBAdminUser{
		Username: "17mei",
		Password: "wothing1708",
	}

	conf.Server.ReportingDisabled = true
	conf.Server.Admin.BindAddress = ":8083"
	conf.Server.BindAddress = ":8088"

	store, err := appdash.NewInfluxDBStore(conf)
	if err != nil {
		log.Fatalf("failed to create influxdb store, error: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Fatalf("fail to close store, error: %v", err)
		}
	}()

	l, err := net.Listen("tcp", ":1726")
	if err != nil {
		log.Fatalf("fail to listen to port, error: %v", err)
	}

	log.Infof("AppDash collector listening on :1726")
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	cs.Debug = true
	cs.Trace = true

	go cs.Start()

	url, err := url.Parse("http://localhost:8700")
	if err != nil {
		log.Fatalf("failed to parsing url, error: %v", err)
	}
	tapp, err := traceapp.New(nil, url)
	if err != nil {
		log.Fatalf("failed to create trace app, error: %v", err)
	}
	tapp.Store = store
	tapp.Queryer = store
	tapp.Aggregator = store
	log.Infof("AppDash web UI running on HTTP :8700")
	log.Fatal(http.ListenAndServe(":8700", tapp))
}
