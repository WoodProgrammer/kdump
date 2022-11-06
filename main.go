package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ackItem = make(map[uint32]SequenceMap)

type SequenceMap struct {
	DateTime      int64
	Count         int
	AckNumber     uint32
	NextSeqNumber uint32
}

type MetricMap struct {
	ipLayer *layers.IPv4
	tcp     *layers.TCP
}

func TcpStream() {

	fmt.Println(reflect.TypeOf(ackItem))
	bpffilter := "tcp"
	handle, err := pcap.OpenLive("en0", 2048, false, -1*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	defer handle.Close()
	handle.SetBPFFilter(bpffilter)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		ipLayer := packet.Layer(layers.LayerTypeIPv4)
		if ipLayer == nil {
			continue
		}

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}

		tcp, ok := tcpLayer.(*layers.TCP)
		ip, ok := ipLayer.(*layers.IPv4)

		metricData := MetricMap{ip, tcp}

		if !ok {
			continue
		}

		if len(tcp.Payload) < 1 {
			continue
		}

		tcpchan <- &metricData
	}

}

func ExportPrometheus() {

	port := flag.String("port", "7070", "display colorized output")
	flag.Parse()

	log.Println("Starting metric server:", *port)

	http.Handle("/metrics", promhttp.Handler())
	portVal := ":" + *port
	http.ListenAndServe(portVal, nil)

}

func main() {
	go RetransmissionHandler()
	go ExportPrometheus()
	TcpStream()
}
