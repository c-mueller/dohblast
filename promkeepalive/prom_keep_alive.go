package promkeepalive

import (
	"fmt"
	"github.com/c-mueller/dohblast/doh"
	"github.com/c-mueller/dohblast/qnamegen"
	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	summary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:  "dohblast",
		Name:       "exec_times",
		Help:       "DoH Response Time",
		Objectives: map[float64]float64{0.1: 0.09, 0.2: 0.08, 0.5: 0.05, 0.8: 0.02, 0.9: 0.01, 0.99: 0.001},
	}, []string{"endpoint", "success"})
	totalCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "dohblast",
		Name:      "total_request_count",
		Help:      "Total Request Count per Endpoint",
	}, []string{"endpoint"})
	failCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "dohblast",
		Name:      "failed_request_count",
		Help:      "Total Request Count per Endpoint",
	}, []string{"endpoint"})
)

func PromKeepAlive(serverEndpoint string, endpoints []string, interval time.Duration, verbose bool) {
	fmt.Println("Initializing Prometheus Metrics...")
	prometheus.MustRegister(summary, totalCount, failCount)

	tickers := make([]time.Ticker, 0)

	for _, v := range endpoints {
		fmt.Printf("Launching Goroutine for %q\n", v)
		t := time.NewTicker(interval)
		go listenerLoop(t, v, verbose)
		tickers = append(tickers, *t)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			for _, v := range tickers {
				v.Stop()
			}
			os.Exit(0)
		}
	}()

	fmt.Printf("Listening on %s\n", serverEndpoint)
	panic(http.ListenAndServe(serverEndpoint, promhttp.Handler()))
}

func listenerLoop(t *time.Ticker, v string, verbose bool) {
	fmt.Println(v)
	for {
		select {
		case <-t.C:
			qn := qnamegen.GenerateRandomQName()
			start := time.Now()

			m := new(dns.Msg)
			m.SetQuestion(qn, dns.TypeA)

			_, err := doh.QueryDoH(v, *m)
			summary.WithLabelValues(v, fmt.Sprintf("%v", err != nil)).Observe(float64(time.Now().Sub(start).Nanoseconds()) / 1000000)
			totalCount.WithLabelValues(v).Add(1)
			if err != nil {
				failCount.WithLabelValues(v).Add(1)
				fmt.Printf("Request has Failed with error: %q on endpoint %q\n", err.Error(), v)
			}
			if verbose {
				fmt.Printf("[%s] %s: Sent A Query to %q in %s\n", time.Now().Format("15:04:05"), v, qn, time.Now().Sub(start).String())
			}
		}
	}
}
