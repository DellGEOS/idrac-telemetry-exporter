package redfishmetricreport

import (
	"log"
	"crypto/tls"
	"net/http"
	"encoding/json"
	"strconv"
	"errors"
	"strings"
//	"sync"

	"github.com/prometheus/client_golang/prometheus"
//	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	username = "username"
	password = "password"
)


var collectors map[string]map[string]*prometheus.GaugeVec

func addGauge(target string, metricValue MetricValue, reportName string, serviceTag string,  registry *prometheus.Registry) {
	var gauge *prometheus.GaugeVec

	log.Printf("%s:\taddGauge:\tCalled with reportName = %s, serviceTag %s and metricValue = %v", target, reportName, serviceTag, metricValue)

	if collectors[target] == nil {
//		log.Printf("%s:\taddGauge:\tcollector for target %s doesn't exist, creating it", target, target)
		collectors[target] = make(map[string]*prometheus.GaugeVec)
	}
//	log.Printf("%s:\taddGauge:\tChecking to see if key %s exists in collectors map %+v", target, metricValue.MetricId, collectors[target])
//	contextID := strings.Replace(metricValue.Oem.Dell.ContextID, " ", "", -1)
	_, keyExists := collectors[target][metricValue.MetricId]
	if keyExists == false {
//		log.Printf("%s:\taddGauge:\tNo entry in collector for metricId %s, creating gauge", target, metricValue.MetricId)
		if strings.Contains(metricValue.Oem.Dell.ContextID, ".") {
			gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: "PowerEdge",
				Subsystem: reportName,
				Name: metricValue.MetricId,
			},
			[]string{"Target", "ServiceTag", "Metric", "FQDD"})
		} else {
			gauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
				Namespace: "PowerEdge",
				Subsystem: reportName,
				Name: metricValue.MetricId,
			},
			[]string{"Target", "ServiceTag", "Metric", "FQDD"})
		}
//		log.Printf("%s:\taddGauge:\tCreated gauge %v for metric %s, adding it to registry", target, *gauge, metricValue.MetricId)
		registry.MustRegister(gauge)
		floatVal, err := strconv.ParseFloat(metricValue.Value, 64)
		if err != nil {
			if metricValue.Value == "Up" || metricValue.Value == "Operational" {
				floatVal = 1
			}
		}
		log.Printf("%s:\taddGauge:\tSetting value for serviceTag %s, with FQDD %s, metric %s to %.2f", target, serviceTag, metricValue.Oem.Dell.FQDD, metricValue.MetricId, floatVal)
		gauge.WithLabelValues(target, serviceTag, metricValue.MetricId, metricValue.Oem.Dell.FQDD).Set(floatVal)
		collectors[target][metricValue.MetricId] = gauge
		//	}
//	if collectors["FQDD"][metricValue.MetricId] == nil {
	} else {
//		log.Printf("%s:\taddGauge:\tKey %s already exists, adding new metric to gauge", target, metricValue.MetricId)
		gauge := collectors[target][metricValue.MetricId]
		floatVal, _ := strconv.ParseFloat(metricValue.Value, 64)
		log.Printf("%s:\taddGauge:\tSetting value for serviceTag %s, with FQDD %s, metric %s to %.2f", target, serviceTag, metricValue.Oem.Dell.FQDD, metricValue.MetricId, floatVal)
		gauge.WithLabelValues(target, serviceTag, metricValue.MetricId, metricValue.Oem.Dell.FQDD).Set(floatVal)
	}
}

func getMetricReport(target string, reportURL string) MetricReport {
	var metricReport MetricReport
	
	log.Printf("%s:\tgetMetricReport:\tCalled getMetricReport for report %s", target, reportURL)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	uri := "https://" + target + reportURL
	req, err := http.NewRequest("GET", uri, nil)
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")
//	log.Printf("%s:\tgetMetricReport:\tGetting report %s from target %s", target, report, target)
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("%s:\tgetMetricReport:\tStatus code = %d, exiting...", response.StatusCode)
		panic(errors.New(response.Status))
	}

//	log.Printf("%s:\tgetMetricReport:\tGot report %s, decoding it", target, report)
	err = json.NewDecoder(response.Body).Decode(&metricReport)
	if err != nil {
		log.Fatal(err)
	}

	response.Body.Close()

	return metricReport
}

func getReportList(target string) []string {
	var metricReportList MetricReportList
	var reports []string = make([]string, 0)

		transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	uri := "https://" + target + "/redfish/v1/TelemetryService/MetricReports"
	req, err := http.NewRequest("GET", uri, nil)
	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")
//	log.Printf("%s:\tgetMetricReport:\tGetting report %s from target %s", target, report, target)
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		log.Printf("%s:\tgetMetricReport:\tStatus code = %d, exiting...", response.StatusCode)
		panic(errors.New(response.Status))
	}

//	log.Printf("%s:\tgetMetricReport:\tGot report %s, decoding it", target, report)
	err = json.NewDecoder(response.Body).Decode(&metricReportList)
	if err != nil {
		log.Fatal(err)
	}

	response.Body.Close()

	for _, member := range metricReportList.Members {
		reports = append(reports, member.ODataId)
	}

	return reports
}

func Probe(target string, registry *prometheus.Registry) bool {
	metricReports := make(map[string][]MetricReport)
	collectors = make(map[string]map[string]*prometheus.GaugeVec)
//	log.Printf("%s:\tredfishmetricreport.Probe:\tGetting reports from target %s", target, target)

	reports := getReportList(target)
	for _, report := range reports {
//		log.Printf("%s:\tredfishmetricreport.Probe:\tGetting report %s", target, report)
		metricReport := getMetricReport(target, report)
		_, keyExists := metricReports[target]
		if keyExists == false {
			metricReports[target] = make([]MetricReport, 0)
		}
		metricReports[target] = append(metricReports[target], metricReport)
	}

	for _, report := range metricReports[target] {
//		log.Printf("%s:\tredfishmetricreport.Probe:\tReceived %d entries from report for target %s", target, report.MetricValuesCount, target)
		for _, metricValue := range report.MetricValues {
//			log.Printf("%s:\tredfishmetricreport.Probe:\tIndex %d, evaluating metricValue = %v", target, idx, metricValue)
			addGauge(target, metricValue, report.Id, report.OemSection.Dell.ServiceTag, registry)
		}
	}
	return true
}

