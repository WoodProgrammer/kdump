package main

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

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

func main() {

	ackItem := make(map[uint32]SequenceMap)
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

		if tcp.Ack != 0 {

			if _, ok := ackItem[tcp.Ack]; ok {
				payloadLength := uint32(len(tcp.Payload))

				if tcp.Seq != ackItem[tcp.Ack].NextSeqNumber {
					fmt.Println("#######START#######")
					fmt.Println("Next SeqNumber:", ackItem[tcp.Ack].NextSeqNumber)
					fmt.Println("Seq:", tcp.Seq)
					fmt.Println("Ack:", tcp.Ack)
					fmt.Println("Dict:", ackItem[tcp.Ack])
					fmt.Println("DataPayload:", uint32(len(tcp.Payload)))
					fmt.Println("#######END#######")
				} else {

				}

				ackPointer := ackItem[tcp.Ack]
				ackPointer.Count = ackPointer.Count + 1
				ackPointer.AckNumber = tcp.Ack

				ackPointer.NextSeqNumber = tcp.Seq + payloadLength

				ackItem[tcp.Ack] = ackPointer

			} else {
				payloadLength := uint32(len(tcp.Payload))
				nextSeqNumber := tcp.Seq + payloadLength

				ackItem[tcp.Ack] = SequenceMap{time.Now().UnixNano(), 0, tcp.Seq, nextSeqNumber}

			}

		} else {
			continue
		}
	}

}
