package main

import (
	_ "encoding/json"
	_ "flag"
	"fmt"
	_ "log"
	"net"
	_ "net/http"
	"os"
	"time"

	_ "github.com/prometheus/client_golang/prometheus"
	_ "github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tatsushid/go-fastping"
)

type PingResult struct {
	destination string `json:"destination"`
	result      int    `json:"rtt"`
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>exping Exporter</title>
		</head>
		<body>
			<h1>exping Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
		</body>
		</html>
	`))
}

func main() {
	p := fastping.NewPinger()
	ra, err := net.ResolveIPAddr("ip4:icmp", os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		pong := new(PingResult)
		pong.destination = addr.String()
		pong.result = 0
		pong_string := fmt.Sprintf("{destination:%s,result:%d}", pong.destination, pong.result)
		fmt.Println(pong_string)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}

}
