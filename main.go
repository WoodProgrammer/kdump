package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

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
	ipLayer    *layers.IPv4
	tcp        *layers.TCP
	deviceName string
}

func TcpStream(device string, filter string) {
	var devicePrefix string

	bpffilter := filter
	rootDeviceList, err := pcap.FindAllDevs()

	if err != nil {
		panic(err)
	}

	if device == "" {
		devicePrefix = "any"

	} else {
		devicePrefix = device
	}

	log.Println("Device to listen is ", devicePrefix)
	handle, err := pcap.OpenLive(devicePrefix, 2048, true, pcap.BlockForever)
	if err != nil {
		fmt.Println(err)
	}
	defer handle.Close()
	handle.SetBPFFilter(bpffilter)

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		index := packet.Metadata().InterfaceIndex

		detectedDeviceName := rootDeviceList[index].Name

		ipLayer := packet.Layer(layers.LayerTypeIPv4)

		packetCount.WithLabelValues(detectedDeviceName).Add(1.0) // count of the package that we received

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
			dfMetric.WithLabelValues(ip.SrcIP.String(), ip.DstIP.String(), detectedDeviceName).Add(1.0)
		}

		for _, val := range tcp.Options {

			if val.OptionType.String() == "WindowScale" {
				windowscaleMetricCount.WithLabelValues(ip.SrcIP.String(), ip.DstIP.String(), detectedDeviceName).Add(1.0)
			}
		}

		if tcp.RST == true {
			rstMetric.WithLabelValues(ip.SrcIP.String(), ip.DstIP.String(), detectedDeviceName).Add(1.0)

		}

		metricData := MetricMap{ip, tcp, detectedDeviceName}

		if !ok {
			continue
		}

		if len(tcp.Payload) < 1 {
			continue
		}

		ipPacketSize.WithLabelValues(detectedDeviceName).Add(float64(len(ip.Payload)))

		tcpPacketSize.WithLabelValues(detectedDeviceName).Add(float64(len(tcp.Payload)))

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
	device := flag.String("device", "", "interface identifier to sniff")
	filter := flag.String("filter", "tcp", "Definition of filter to run")

	flag.Parse()

	go RetransmissionHandler()
	go ExportPrometheus(*port, *filter)
	TcpStream(*device, *filter)
}
