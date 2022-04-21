FROM alpine:latest

RUN apk add --no-cache smartmontools

ADD bin/smartctl_exporter /bin/smartctl_exporter

EXPOSE 9111

ENTRYPOINT "/bin/smartctl_exporter"
