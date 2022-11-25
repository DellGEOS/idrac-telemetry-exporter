# iDRAC Telemetry Prometheus Exporter

iDRAC is the Baseboard Management Computer (BMC) embedded in all Dell PowerEdge servers. The iDRAC offers a telemetry streaming service to export and expose numerous metrics pertaining to the health and performance of the server.

Prometheus is the de-facto standard for collecting metrics within Kubernetes environments. Prometheus works on the basis of scraping data from exporters. 
This exporter collects metrics from the iDRAC telemetry streaming service and exposes them to Prometheus to scrape. It uses Redfish endpoint to gather the metrics and gathers metrics from all the enabled metric report.

## Running the exporter
The easiest way to run the iDRAC Telemetry exporter is by leveraging the docker container for this exporter. 
```
docker run -d --name idrac-telemetry-exporter -p 3355:3355 dbblackdiamond/idrac-telemetry-exporter:latest
```

## Configuration
To configure Prometheus to scrape the exporter endpoint, the following lines will need to be added to the Prometheus configuration file.
```
  - job_name: 'idrac_telemetry'
    scrape_interval: 2m
    scrape_timeout: 45s
    metrics_path: /probe
    static_configs:
      - targets:
        - <idrac ip address>
        - <idrac ip address>
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: <exporter ip address>:3355
```

## License

Apache 2.0
