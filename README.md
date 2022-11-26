# iDRAC Telemetry Prometheus Exporter

iDRAC is the Baseboard Management Computer (BMC) embedded in all Dell PowerEdge servers. The iDRAC offers a telemetry streaming service to export and expose numerous metrics pertaining to the health and performance of the server.

Prometheus is the de-facto standard for collecting metrics within Kubernetes environments. Prometheus works on the basis of scraping data from exporters. 
This exporter collects metrics from the iDRAC telemetry streaming service and exposes them to Prometheus to scrape. It uses Redfish endpoint to gather the metrics and gathers metrics from all the enabled metric report.

## Running the exporter
The easiest way to run the iDRAC Telemetry exporter is by leveraging the docker container for this exporter. 
```
docker run -d --name idrac-telemetry-exporter -p 3355:3355 -v <configdirector>/config.yml:/etc/idrac-telemetry-exporter/config.yml dbblackdiamond/idrac-telemetry-exporter:latest
```
By default, the exporter will listen to port 3355.
The `config.yml` file contains the configuration for the exporter, primarily, the `username` and `password` to authenticate against each of the targets.

## Exporter Configuration
The configuration for the exporter is included in the `config.yml` file provided in argument to the `docker run` command. It contains 2 sections: `global` and `idracs`. The `global` section can be used to specify global usernames and passwords. These will be used if not specific iDRAC's username and password are defined. 
The specific iDRAC's username and password are defined in the `idracs` section. The section also includes the `ip address` of the iDRAC and need to match the IP address specified in the `target` section of the Prometheus scraper configuration.

```
gobal:
  username: "<username>"
  password: "<password>"
  
idracs:
  - <hostname 1>:
    address: "<ip address>"
    username: "<username>"
    password: "<password>"
  - <hostname 2>
    address: "<ip address>"
    username: "<username>"
    password: "<password>"
```

## Prometheus Scraper Configuration
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
Because the exporter is pulling all the redfish reports for each scrape, it is recommended not to use a scrape interval lower than 2 minutes. Pulling all the redfish reports and processing them takes around 30 seconds, hence the recommended timeout value of 45 seconds.

## License

Apache 2.0
