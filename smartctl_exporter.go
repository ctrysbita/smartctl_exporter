package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	devices    []*Device
	results    = sync.Map{}
	collectors = map[string]*prometheus.GaugeVec{}
)

func WithMetrics(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		wg := sync.WaitGroup{}
		wg.Add(len(devices))
		for _, d := range devices {
			go func(d *Device) {
				results.Store(d.Name, GetAll(d))
				wg.Done()
			}(d)
		}
		wg.Wait()

		for _, d := range devices {
			if r, ok := results.Load(d.Name); ok {
				r := r.(*Result)
				passed := 0
				if r.Passed {
					passed = 1
				}
				collectors["device_status"].
					WithLabelValues(d.Name, r.ModelName, r.SerialNumber, r.FirmwareVersion).
					Set(float64(passed))

				for k, v := range r.Attributes {
					c, ok := collectors[k]
					if !ok {
						c = promauto.NewGaugeVec(
							prometheus.GaugeOpts{
								Namespace: "smartctl",
								Name:      k,
							},
							[]string{"device"},
						)
						collectors[k] = c
					}
					c.WithLabelValues(d.Name).Set(v)
				}
			}
		}

		handler.ServeHTTP(w, r)
	}
}

func main() {
	devices = GetDevices()

	deviceStatusCollector := promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "smartctl",
			Name:      "device_status",
			Help:      "Device Status",
		},
		[]string{
			"device",
			"model_name",
			"serial_number",
			"firmware_version",
		},
	)
	collectors["device_status"] = deviceStatusCollector

	promHandler := promhttp.Handler()
	http.HandleFunc("/scrape", WithMetrics(promHandler))
	log.Fatal(http.ListenAndServe(":9111", nil))
}
