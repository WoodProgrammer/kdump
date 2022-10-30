package main

import "fmt"

type SequenceMap struct {
	Count uint32
}

func incrementCount(s *SequenceMap) {
	s.Count = s.Count + 1
}

func main() {
	ackItem := SequenceMap{}
	ackItem.Count = 123
	incrementCount(&ackItem)

	fmt.Println(ackItem.Count)
}
