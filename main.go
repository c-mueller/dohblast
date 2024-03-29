package main

import (
	"encoding/json"
	"fmt"
	"github.com/c-mueller/dohblast/blast"
	"github.com/c-mueller/dohblast/keepalive"
	"github.com/c-mueller/dohblast/promkeepalive"
	"github.com/c-mueller/dohblast/qnamegen"
	"gopkg.in/alecthomas/kingpin.v2"
	"math/rand"
	"time"
)

var (
	blastCmd = kingpin.Command("blast", "Send Random A Requests to the given Endooint").Alias("b")

	threadCountBlast = blastCmd.Command("threads", "Send requests by count of background daemons").Alias("t")
	endpoint         = threadCountBlast.Arg("endpoint", "Endpoint").Required().String()
	threadCount      = threadCountBlast.Arg("thread-count", "Thread Count").Required().Int()

	rateBlast = blastCmd.Command("rate", "Send requests by count per second").Alias("r")
	rendpoint = rateBlast.Arg("endpoint", "Endpoint").Required().String()
	rate      = rateBlast.Arg("rate", "The number of requests per second to send").Required().Int()

	tldWeightsCmd = kingpin.Command("tld-weights", "output the default tld weights as json")

	keepAliveCmd = kingpin.Command("keep-alive", "Send Requests in a defined interfval")
	kaEndpoint   = keepAliveCmd.Arg("endpoint", "Endpoint").Required().String()
	kaDuration   = keepAliveCmd.Arg("interval", "Interval").Default("5s").Duration()

	promKeepAlive     = kingpin.Command("prom-keep-alive", "Run Keepalive as a server")
	pkaServerEndpoint = promKeepAlive.Arg("server-endpoint", "Server Endpoint").Default(":9898").String()
	pkaEndpoints      = promKeepAlive.Flag("endpoint", "doh endpoint").Short('e').Required().Strings()
	pkaDuration       = promKeepAlive.Arg("interval", "Interval").Default("5s").Duration()

	verbose = kingpin.Flag("verbose", "Verbose output").Short('v').Default("false").Bool()

	cmd string
)

func init() {
	cmd = kingpin.Parse()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	switch cmd {
	case "blast thread-count":
		blast.BlastByThreadCountCommand(*endpoint, *threadCount)
		break
	case "blast rate":
		blast.BlastByRateCommand(*rendpoint, *rate)
		break
	case "keep-alive":
		keepalive.KeepAliveCmd(*kaEndpoint, *kaDuration, *verbose)
		break
	case "prom-keep-alive":
		promkeepalive.PromKeepAlive(*pkaServerEndpoint, *pkaEndpoints, *pkaDuration, *verbose)
		break
	case "tld-weights":
		data, _ := json.Marshal(qnamegen.DefaultTLDList)
		fmt.Println(string(data))
		break
	}
}
