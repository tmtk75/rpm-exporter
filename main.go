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

type myCollector struct {
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

func (c myCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- desc
}

func (c myCollector) Collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(1), "foo", "1.2.3")
	ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, float64(1), "bar", "0.1.2")
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

var desc *prometheus.Desc

func main() {
	flag.Var(&nameFlags, "name", "rpm name. multiple.")
	flag.Parse()

	desc = prometheus.NewDesc("rpm_info", "Show RPM info", []string{"rpm_name", "version"}, nil)

	var c myCollector
	prometheus.MustRegister(c)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(*addr, nil))
}
