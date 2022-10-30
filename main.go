package main

import (
	"fmt"
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

func (s SequenceMap) setAckItem() {
	s.Count = 0
	s.DateTime = time.Now().UnixNano()
}

func Retransmission() {

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
		if !ok {
			continue
		}

		if len(tcp.Payload) < 1 {
			continue
		}
		tcpchan <- tcp

	}

}

func ExportPrometheus() {
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9000", nil)
}

func main() {

	go TcpMetricHandler()
	Retransmission()
}
