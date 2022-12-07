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
	WindowSize    uint16
}

type MetricMap struct {
	ipLayer *layers.IPv4
	tcp     *layers.TCP
}

func TcpStream(ninterface string, filter string) {

	fmt.Println(reflect.TypeOf(ackItem))
	bpffilter := filter
	handle, err := pcap.OpenLive(ninterface, 2048, false, -1*time.Second)
	if err != nil {
		fmt.Println(err)
	}
	defer handle.Close()
	handle.SetBPFFilter(bpffilter)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		fmt.Println("Packet ", packet.Data())

		ipLayer := packet.Layer(layers.LayerTypeIPv4)

		packetCount.WithLabelValues().Add(1.0) // count of the package that we received

		if ipLayer == nil {
			continue
		}

		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue
		}

		tcp, ok := tcpLayer.(*layers.TCP)

		ip, ok := ipLayer.(*layers.IPv4)

		if ip.Flags.String() == "DF" {
			dfMetric.WithLabelValues(ip.SrcIP.String(), ip.DstIP.String()).Add(1.0)
		}

		for _, val := range tcp.Options {

			if val.OptionType.String() == "WindowScale" {
				windowscaleMetricCount.WithLabelValues(ip.SrcIP.String(), ip.DstIP.String()).Add(1.0)
			}
		}

		if tcp.RST == true {
			rstMetric.WithLabelValues(ip.SrcIP.String(), ip.DstIP.String()).Add(1.0)

		}

		metricData := MetricMap{ip, tcp}

		if !ok {
			continue
		}

		if len(tcp.Payload) < 1 {
			continue
		}

		ipPacketSize.WithLabelValues().Add(float64(len(ip.Payload)))

		tcpPacketSize.WithLabelValues().Add(float64(len(tcp.Payload)))

		tcpchan <- &metricData

	}

}

func ExportPrometheus(port string, filter string) {

	log.Println("Starting metric server:", port)
	log.Println("Bpf filter to run:", filter)

	http.Handle("/metrics", promhttp.Handler())
	portToRun := ":" + port
	http.ListenAndServe(portToRun, nil)

}

func main() {
	port := flag.String("port", "9090", "port to run")
	ninterface := flag.String("interface", "en0", "interface identifier to sniff")
	filter := flag.String("filter", "tcp", "Definition of filter to run")

	flag.Parse()

	go RetransmissionHandler()
	go ExportPrometheus(*port, *filter)
	TcpStream(*ninterface, *filter)
}
