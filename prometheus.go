package main

import (
	"time"

	"github.com/google/gopacket/layers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	retranmissionMetricCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "retransmission_metric",
		Help: "TCP Retransmission metrics that contains details and error type",
	}, []string{"dstPort", "srcPort"})
)

var tcpchan chan *layers.TCP

func TcpMetricHandler() {
	tcpchan = make(chan *layers.TCP, 100)

	for tcpData := range tcpchan {

		if tcpData.Ack != 0 {

			if _, ok := ackItem[tcpData.Ack]; ok {
				payloadLength := uint32(len(tcpData.Payload))

				if tcpData.Seq != ackItem[tcpData.Ack].NextSeqNumber {
					if ackItem[tcpData.Ack].Count != 0 {
						retranmissionMetricCount.WithLabelValues(tcpData.DstPort.String(), tcpData.SrcPort.String()).Add(float64(ackItem[tcpData.Ack].Count))
					}
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
