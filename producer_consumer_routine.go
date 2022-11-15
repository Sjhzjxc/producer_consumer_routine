package producer_consumer_routine

import (
	"context"
	"sync"
)

type ProducerConsumer struct {
	Wg            *sync.WaitGroup
	Cancel        context.CancelFunc
	Ctx           context.Context
	Chan          chan interface{}
	Producer      func(func(interface{}) bool)
	ProducerCount int
	Consumer      func(interface{})
	ConsumerCount int
	Params        map[string]interface{}
}

func NewProducerConsumer(producer func(func(interface{}) bool), consumer func(interface{}), chanCount, producerCount, consumerCount int, params map[string]interface{}) *ProducerConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProducerConsumer{
		Wg:            &sync.WaitGroup{},
		Cancel:        cancel,
		Ctx:           ctx,
		Chan:          make(chan interface{}, chanCount),
		Producer:      producer,
		ProducerCount: producerCount,
		Consumer:      consumer,
		ConsumerCount: consumerCount,
		Params:        params,
	}
}

func DefaultProducerConsumer(producer func(func(interface{}) bool), consumer func(interface{})) *ProducerConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProducerConsumer{
		Wg:            &sync.WaitGroup{},
		Cancel:        cancel,
		Ctx:           ctx,
		Chan:          make(chan interface{}, 1),
		Producer:      producer,
		ProducerCount: 1,
		Consumer:      consumer,
		ConsumerCount: 1,
		Params:        nil,
	}
}

func (c *ProducerConsumer) producer(worker func(func(interface{}) bool)) {
	defer c.Defer()
	isBreak := false

	worker(func(i interface{}) bool {
		select {
		case <-c.Ctx.Done():
			isBreak = true
		default:
			c.Chan <- i
		}
		return isBreak
	})

	if !isBreak {
		c.SendNil()
	}
}

func (c *ProducerConsumer) consumer(worker func(interface{})) {
	defer c.Defer()
LOOP:
	for {
		select {
		case <-c.Ctx.Done():
			break LOOP
		case value := <-c.Chan:
			if value == nil {
				break LOOP
			}
			worker(value)
		default:
		}
	}
}

func (c *ProducerConsumer) Routine() {
	for i := 0; i < c.ProducerCount; i++ {
		c.Wg.Add(1)
		//fmt.Println("wg.Add 1")
		go c.producer(c.Producer)
	}
	for i := 0; i < c.ConsumerCount; i++ {
		c.Wg.Add(1)
		//fmt.Println("wg.Add 1")
		go c.consumer(c.Consumer)
	}
	//fmt.Println("over111")
	c.Wg.Wait()
}

func (c *ProducerConsumer) Defer() {
	c.Wg.Done()
}

func (c *ProducerConsumer) Close() {
	close(c.Chan)
}

func (c *ProducerConsumer) SendNil() {
	for i := 0; i < c.ConsumerCount; i++ {
		c.Chan <- nil
	}
}
