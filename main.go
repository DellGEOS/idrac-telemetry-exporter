package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"gopkg.in/yaml.v3"

	"idrac_telemetry_exporter/prober"
	"idrac_telemetry_exporter/redfishmetricreport"
)

var configFile = "/etc/idrac-telemetry-exporter/config.yml"

func init() {
	prometheus.MustRegister(version.NewCollector("idrac_telemetry_exporter"))
}

func main() {

	var config redfishmetricreport.Config

	if os.Getenv("CONFIGFILE") != "" {
		configFile = os.Getenv("CONFIGFILE")
	}
	configData, err := os.ReadFile(configFile)
	if err != nil {
		log.Printf("main: Error reading config file %s, err = %v", configFile, err)
		log.Fatal(err)
	}
	err = yaml.Unmarshal([]byte(configData), &config)
	if err != nil {
		log.Printf("main: Error unmarshaling config data, err = %v", err)
		log.Fatal(err)
	}
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", func(w http.ResponseWriter, req *http.Request) {
		prober.Handler(w, req, config, nil)
	})
	log.Println("main: Serving requests on port 3355")
	err = http.ListenAndServe(":3355", nil)
	if err != nil {
		log.Printf("main: Failed to start webserver %v", err)
		log.Fatal(err)
	}
}
