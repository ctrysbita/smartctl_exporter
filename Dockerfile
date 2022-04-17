FROM alpine:latest

RUN apk update; \
    apk add --no-cache python3 py3-pip smartmontools; \
    python3 -m pip install prometheus_client;

ADD smartctl_exporter smartctl_exporter

EXPOSE 9111
ENTRYPOINT "smartctl_exporter/smartctl_exporter.py"