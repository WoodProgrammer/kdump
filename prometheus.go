package main

import (
	"fmt"
	"time"

	"github.com/google/gopacket/layers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	tcpMetricCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tcp_metric",
		Help: "TCP metrics that contains details and error type",
	}, []string{"dstPort", "srcPort", "error", "srcHost"})
)

var tcpchan chan *layers.TCP

func TcpMetricHandler() {
	tcpchan = make(chan *layers.TCP, 100)

	for tcpData := range tcpchan {

		if tcpData.Ack != 0 {

			if _, ok := ackItem[tcpData.Ack]; ok {
				payloadLength := uint32(len(tcpData.Payload))

				if tcpData.Seq != ackItem[tcpData.Ack].NextSeqNumber {
					fmt.Println("#######START#######")
					fmt.Println("Next SeqNumber:", ackItem[tcpData.Ack].NextSeqNumber)
					fmt.Println("Seq:", tcpData.Seq)
					fmt.Println("Ack:", tcpData.Ack)
					fmt.Println("Dict:", ackItem[tcpData.Ack])
					fmt.Println("DataPayload:", uint32(len(tcpData.Payload)))
					fmt.Println("#######END#######")

				} else {

				}
				ackPointer := ackItem[tcpData.Ack]
				ackPointer.Count = ackPointer.Count + 1
				ackPointer.AckNumber = tcpData.Ack

				ackPointer.NextSeqNumber = tcpData.Seq + payloadLength

				ackItem[tcpData.Ack] = ackPointer

			} else {
				payloadLength := uint32(len(tcpData.Payload))
				nextSeqNumber := tcpData.Seq + payloadLength
				ackItem[tcpData.Ack] = SequenceMap{time.Now().UnixNano(), 0, tcpData.Seq, nextSeqNumber}
			}

		} else {
			continue
		}

	}
}
