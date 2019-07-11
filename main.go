package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	namespace = "rpm"
)

type myCollector struct{
	gauges []prometheus.Gauge
}

func gv(name string) string {
	b := bytes.NewBuffer([]byte{})
	ctx := context.Background()
	cmd := exec.CommandContext(ctx, "/usr/bin/rpm", "-q", name)
	cmd.Stderr = b
	cmd.Stdout = b
	if err := cmd.Run(); err != nil {
		fmt.Println(string(b.Bytes()))
		log.Fatal(err)
	}
	return strings.TrimSpace(strings.TrimLeft(string(b.Bytes()), name+"-"))
}

func newGauge(name string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        "info",
		Help:        "Info of " + name,
		ConstLabels: map[string]string{"version": gv(name), "name": name},
	})
}

func (c myCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, g := range c.gauges {
		ch <- g.Desc()
	}
}

func (c myCollector) Collect(ch chan<- prometheus.Metric) {
	for _, g := range c.gauges {
		ch <- prometheus.MustNewConstMetric(g.Desc(), prometheus.GaugeValue, float64(1))
	}
}

var addr = flag.String("listen-address", "0.0.0.0:9872", "The address to listen on for HTTP requests.")

type stringFlags []string

func (i *stringFlags) String() string {
	return "string flags"
}

func (i *stringFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var nameFlags stringFlags

func main() {
	flag.Var(&nameFlags, "name", "rpm name. multiple.")
	flag.Parse()

	var c myCollector
	// c.gauges = append(c.gauges, newGauge("glibc"))
	c.gauges = append(c.gauges, newGauge("openssl"))

	prometheus.MustRegister(c)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
