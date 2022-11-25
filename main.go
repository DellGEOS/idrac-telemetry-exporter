package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"

//	"github.com/gorilla/mux"
//	"idrac_telemetry_exporter/redfish"
	"idrac_telemetry_exporter/prober"
)

//var (
//	CPUMetrics []CpuMemCollector
//)

//func handler() {
//	GetRedfishReport("/redfish/v1/TelemetryService/MetricReports/CPUMemMetrics")
//	ConvertRedfishToPrometheus()

//}

func init() {
	prometheus.MustRegister(version.NewCollector("idrac_telemetry_exporter"))
}

func main() {


//	metricReports := redfish.GetRedFishReports()
//	log.Printf("main: Received %d reports from redfish.", len(metricReports))
//	log.Printf("main: Decoding them")
//	registry := prometheus.NewRegistry()
//	log.Printf("main: creating new collector for CPUMem metrics")
//	cpuMemCollector := cpumemcollector.NewCpuMemCollector()
//	log.Printf("main: adding collector %v to registry", cpuMemCollector)
//	registry.MustRegister(cpuMemCollector)
//	log.Printf("main: collector added to registry")

//	go getCpuMemMetrics(registry)
//	gatherer := prometheus.Gatherer(registry)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/probe", func(w http.ResponseWriter, req *http.Request) {
		prober.Handler(w, req, nil)
	})
	log.Println("Serving requests on port 3355")
	err := http.ListenAndServe(":3355", nil)
	if err != nil {
		log.Printf("main: Failed to start webserver %v", err)
		log.Fatal(err)
	}
}
