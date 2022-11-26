package producer_consumer_routine

import (
	"context"
	"sync"
)

type ProducerConsumerContent struct {
	Producer      func(func(interface{}) bool)
	ProducerCount int
	Consumer      func(interface{})
	ConsumerCount int
}

type ProducerConsumer struct {
	wg       *sync.WaitGroup
	Wg       *sync.WaitGroup
	Cancel   context.CancelFunc
	Ctx      context.Context
	Chan     chan interface{}
	Contents []ProducerConsumerContent
	Params   map[string]interface{}
	Status   int
}

func NewProducerConsumer(contents []ProducerConsumerContent, chanCount int, params map[string]interface{}) *ProducerConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProducerConsumer{
		wg:       &sync.WaitGroup{},
		Wg:       &sync.WaitGroup{},
		Cancel:   cancel,
		Ctx:      ctx,
		Chan:     make(chan interface{}, chanCount),
		Contents: contents,
		Params:   params,
	}
}

func DefaultProducerConsumer(contents []ProducerConsumerContent) *ProducerConsumer {
	ctx, cancel := context.WithCancel(context.Background())
	return &ProducerConsumer{
		wg:       &sync.WaitGroup{},
		Wg:       &sync.WaitGroup{},
		Cancel:   cancel,
		Ctx:      ctx,
		Chan:     make(chan interface{}, 1),
		Contents: contents,
		Params:   nil,
	}
}

func (c *ProducerConsumer) producer(worker func(func(interface{}) bool), count int) {
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
		c.SendNil(count)
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
	if c.Ctx.Err() != nil {
		ctx, cancel := context.WithCancel(context.Background())
		c.Ctx = ctx
		c.Cancel = cancel
	}
	c.Wg.Add(len(c.Contents))
	for _, content := range c.Contents {
		c.wg.Add(content.ConsumerCount + content.ProducerCount)
		for i := 0; i < content.ProducerCount; i++ {
			go c.producer(content.Producer, content.ConsumerCount)
		}
		for i := 0; i < content.ConsumerCount; i++ {
			go c.consumer(content.Consumer)
		}
		c.wg.Wait()
		c.Wg.Done()
	}
	c.Wg.Wait()
}

func (c *ProducerConsumer) Defer() {
	c.wg.Done()
}

func (c *ProducerConsumer) Close() {
	if c.Chan != nil {
		close(c.Chan)
		c.Chan = nil
	}
}

func (c *ProducerConsumer) SetChanCount(count int) {
	c.Close()
	c.Chan = make(chan interface{}, count)
}

func (c *ProducerConsumer) SendNil(count int) {
	for i := 0; i < count; i++ {
		c.Chan <- nil
	}
}
