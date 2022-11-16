package main

import (
	"fmt"
	"github.com/Sjhzjxc/producer_consumer_routine"
	"time"
)

func producer(executor func(interface{}) bool) {
LOOP:
	for i := 0; i < 10; i++ {
		isBreak := executor(i)
		if isBreak {
			break LOOP
		}
		time.Sleep(time.Second)
	}
}

func producer2(executor func(interface{}) bool) {
LOOP:
	for i := 0; i < 10; i++ {
		isBreak := executor(i + 100)
		if isBreak {
			break LOOP
		}
		time.Sleep(time.Second)
	}
}

func consumer(value interface{}) {
	fmt.Println(value)
}

func consumer2(value interface{}) {
	fmt.Println(value, value.(int)+1)
}

func main() {
	content := producer_consumer_routine.ProducerConsumerContent{
		Producer:      producer,
		ProducerCount: 1,
		Consumer:      consumer,
		ConsumerCount: 1,
	}
	content2 := producer_consumer_routine.ProducerConsumerContent{
		Producer:      producer2,
		ProducerCount: 1,
		Consumer:      consumer2,
		ConsumerCount: 1,
	}
	contents := []producer_consumer_routine.ProducerConsumerContent{
		content,
		content2,
	}
	pc := producer_consumer_routine.DefaultProducerConsumer(contents)
	go pc.Routine()
	time.Sleep(10 * time.Second)
	pc.Cancel()
	pc.Wg.Wait()
	go pc.Routine()
	time.Sleep(10 * time.Second)
	pc.Cancel()
	pc.Wg.Wait()
}
