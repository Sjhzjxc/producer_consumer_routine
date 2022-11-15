package main

import (
	"fmt"
	"github.com/Sjhzjxc/producer_consumer_routine"
	"time"
)

func producer(executor func(interface{}) bool) {
LOOP:
	for i := 0; i < 1000; i++ {
		isBreak := executor(i)
		if isBreak {
			break LOOP
		}
		time.Sleep(time.Second)
	}
}

func consumer(value interface{}) {
	fmt.Println(value)
}

func main() {
	pc := producer_consumer_routine.DefaultProducerConsumer(producer, consumer)
	go pc.Routine()
	time.Sleep(10 * time.Second)
	pc.Cancel()
	pc.Wg.Wait()
}
