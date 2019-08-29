package blast

import (
	"fmt"
	"github.com/c-mueller/dohblast/doh"
	"github.com/c-mueller/dohblast/qnamegen"
	"github.com/miekg/dns"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

func BlastCommand(endpoint string, threads int) {
	var ops uint64
	var failCnt uint64

	var wg sync.WaitGroup
	chns := make([]chan interface{}, 0)

	for i := 0; i < threads; i++ {
		chn := make(chan interface{})
		wg.Add(1)
		go func() {
			for {
				select {
				default:
					qname := qnamegen.GenerateRandomQName()
					m := new(dns.Msg)
					m.SetQuestion(qname, dns.TypeA)

					_, err := doh.QueryDoH(endpoint, *m)
					if err != nil {
						atomic.AddUint64(&failCnt, 1)
					}

					atomic.AddUint64(&ops, 1)
				case <-chn:
					wg.Done()
					return
				}
			}
		}()
		chns = append(chns, chn)
	}

	killchn := make(chan interface{})

	start := time.Now()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			fmt.Println()
			for _, v := range chns {
				close(v)
			}
			close(killchn)
			diff := time.Now().Sub(start)
			fmt.Printf("Sent %d Requests in %s\n", ops, diff.String())
			os.Exit(0)
		}
	}()

	for {
		select {
		default:
			duration := time.Now().Sub(start)
			failRate := (float64(failCnt) / float64(ops)) * 100
			fmt.Printf("\rPerformed %d Requests (%d R/s). %d have Failed. Failure Rate: %f%%", ops, ops/uint64(duration.Seconds()+1), failCnt, failRate)
			time.Sleep(time.Second)
		case <-killchn:
			fmt.Println()
			fmt.Printf("Stopped after %s\n", time.Now().Sub(start))
		}
	}
}
