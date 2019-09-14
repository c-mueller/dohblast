package keepalive

import (
	"fmt"
	"github.com/c-mueller/dohblast/doh"
	"github.com/c-mueller/dohblast/qnamegen"
	"github.com/miekg/dns"
	"os"
	"os/signal"
	"time"
)

func KeepAliveCmd(endpoint string, interval time.Duration, verbose bool) {
	t := time.NewTicker(interval)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			t.Stop()
			os.Exit(0)
		}
	}()

	for {
		select {
		case <-t.C:
			qn := qnamegen.GenerateRandomQName()
			start := time.Now()

			m := new(dns.Msg)
			m.SetQuestion(qn, dns.TypeA)

			res, err := doh.QueryDoH(endpoint, *m)
			if err != nil {
				fmt.Printf("Request has Failed with error: %q\n", err.Error())
			}
			if verbose {
				fmt.Printf("[%s]: Sent A Query to %q in %s\n", time.Now().Format("15:04:05"), qn, time.Now().Sub(start).String())
				fmt.Println(res.String())
			}
		}
	}
}
