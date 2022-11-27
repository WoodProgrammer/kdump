package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	retranmissionMetricCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "retransmission_metric",
		Help: "TCP Retransmission metrics that contains details and error type",
	}, []string{"srcIp", "dstIp"})

	durationMetricCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "duration_metric",
		Help: "TCP Duration time period that contains details and error type",
	}, []string{"srcIp", "dstIp"})

	windowscaleMetricCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "window_scale_metric",
		Help: "WindowScale Detection",
	}, []string{"srcIp", "dstIp"})

	rstMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "rst_metric",
		Help: "TCP RST metrics that contains details and error type",
	}, []string{"srcIp", "dstIp"})

	dfMetric = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ip_df_metric",
		Help: "Ip layer do not fragment metric",
	}, []string{"srcIp", "dstIp"})

	ipPacketSize = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "ip_packet_size",
		Help: "Size of the ip package during timeline",
	}, []string{})

	tcpPacketSize = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "tcp_packet_size",
		Help: "Size of the tcp packet at all",
	}, []string{})

	packageCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "packageCount",
		Help: "Count of the package that received",
	}, []string{})

	zerowindowMetricCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "zerowindow_metric",
		Help: "TCP Zero window metrics that contains details and error type",
	}, []string{"srcIp", "dstIp"})
)

var tcpchan chan *MetricMap

func RetransmissionHandler() {
	tcpchan = make(chan *MetricMap, 5000)

	for metricLaterData := range tcpchan {
		tcpData := metricLaterData.tcp
		ipData := metricLaterData.ipLayer

		if tcpData.Ack != 0 {
			if tcpData.Window == 0 {

				zerowindowMetricCount.WithLabelValues(ipData.SrcIP.String(), ipData.DstIP.String()).Add(1.0)
			}

			currentDateTime := ackItem[tcpData.Seq].DateTime

			duration := time.Nanosecond * time.Duration(time.Now().UnixNano()-currentDateTime)

			if duration != 0 {
				durationMetricCount.WithLabelValues(ipData.SrcIP.String(), ipData.DstIP.String()).Add(float64(duration))
			}

			if _, ok := ackItem[tcpData.Ack]; ok {
				payloadLength := uint32(len(tcpData.Payload))

				if tcpData.Seq != ackItem[tcpData.Ack].NextSeqNumber {
					if ackItem[tcpData.Ack].Count != 0 {
						retranmissionMetricCount.WithLabelValues(ipData.SrcIP.String(), ipData.DstIP.String()).Add(float64(ackItem[tcpData.Ack].Count))

					}

				} else {

				}

				ackPointer := ackItem[tcpData.Ack]
				ackPointer.Count = ackPointer.Count + 1
				ackPointer.AckNumber = tcpData.Ack
				ackPointer.WindowSize = tcpData.Window

				ackPointer.NextSeqNumber = tcpData.Seq + payloadLength

				ackItem[tcpData.Ack] = ackPointer

			} else {
				payloadLength := uint32(len(tcpData.Payload))
				nextSeqNumber := tcpData.Seq + payloadLength
				ackItem[tcpData.Ack] = SequenceMap{time.Now().UnixNano(), 0, tcpData.Seq, nextSeqNumber, tcpData.Window}

			}
		} else {
			continue
		}

	}
}
