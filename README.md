# smartctl_exporter

[![Docker Image Size](https://badgen.net/docker/size/ctrysbita/smartctl_exporter?icon=docker&label=image%20size)](https://hub.docker.com/r/ctrysbita/smartctl_exporter)

Prometheus exporter for [smartmontools](https://www.smartmontools.org/) to export the S.M.A.R.T. attributes.

## Deployment

```sh
docker run --detach --privileged -p 9111:9111 ctrysbita/smartctl_exporter:latest
```

Metrics will be available at http://localhost:9111/scrape

## Grafana Dashboard

![](https://user-images.githubusercontent.com/8357481/164513823-175e1d32-2ba1-41a8-a7f4-76b1b5f07f09.png)
