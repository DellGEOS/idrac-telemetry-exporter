package prober

import (
	"log"
	"net/url"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"idrac_telemetry_exporter/redfishmetricreport"
)

func Handler(w http.ResponseWriter, req *http.Request, config redfishmetricreport.Config, params url.Values) {
	if params == nil {
		params = req.URL.Query()
	}

	probeSuccessGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "probe_success",
		Help: "Display whether or not the probe was a success",
	})
	probeDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "probe_duration_seconds",
		Help: "Returns how long the probe took to complete in seconds",
	})

	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		log.Printf("Prober.Handler:\tTarget parameter missing")
		return
	}

	start := time.Now()
	registry := prometheus.NewRegistry()
	registry.MustRegister(probeSuccessGauge)
	registry.MustRegister(probeDurationGauge)
	success := redfishmetricreport.Probe(target, config, registry)
	duration := time.Since(start).Seconds()
	probeDurationGauge.Set(duration)
	if success {
		probeSuccessGauge.Set(1)
		log.Printf("Prober.Handler:\tProbe succeeded, duration = %.2f", duration)
	} else {
		log.Printf("Prober.Handler:\tProbe failed, duration = %.2f", duration)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, req)
}