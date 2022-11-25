FROM ubuntu:latest

COPY build/idrac-telemetry-exporter /bin/idrac-telemetry-exporter

CMD ["/bin/idrac-telemetry-exporter"]