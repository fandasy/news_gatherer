package eventconsumer

import (
	"context"
	"log"
	"telegramBot/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

var stopSignal bool

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Start() error {

	log.Print("Consumer started")

	defer func() { log.Print("Consumer finished the job") }()

	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			if !stopSignal {
				continue
			} else {
				break
			}
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Print(err)
		}

		if !stopSignal {
			continue
		} else {
			break
		}
	}

	return nil
}

func Stop() {
	stopSignal = true
}

func (c *Consumer) handleEvents(eventsArr []events.Event) error {
	for _, eventsElement := range eventsArr {
		go func(event events.Event) {
			tw := time.Now()
			log.Printf("got new event: %s", event.Text)

			if err := c.processor.Process(context.TODO(), event); err != nil {
				log.Printf("can't handle event: %s", err.Error())
			}

			log.Printf("The %s event was over in %v", event.Text, time.Since(tw))
		}(eventsElement)
	}
	return nil
}
