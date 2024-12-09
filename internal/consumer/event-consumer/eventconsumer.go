package eventconsumer

import (
	"context"
	"log/slog"
	"telegramBot/internal/events"
	"telegramBot/internal/lib/logger/sl"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
	log       *slog.Logger
	stopCh    chan struct{}
}

var stopSignal bool

func New(fetcher events.Fetcher, processor events.Processor, batchSize int, log *slog.Logger) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
		log:       log,
		stopCh:    make(chan struct{}),
	}
}

func (c *Consumer) Start(timeout time.Duration) {

	c.log.Info("Consumer started")

	for {
		select {
		case <-c.stopCh:

			c.log.Info("Consumer stopped")
			
			time.Sleep(5 * time.Second)
			return

		default:

			gotEvents, err := c.fetcher.Fetch(context.TODO(), c.batchSize)
			if err != nil {
				c.log.Error("[ERR] consumer: %s", sl.Err(err))

				continue
			}

			if len(gotEvents) == 0 {
				time.Sleep(timeout)
			}

			c.handleEvents(context.TODO(), gotEvents)

		}
	}
}

func (c *Consumer) Stop() {
	close(c.stopCh)
}

func (c *Consumer) handleEvents(ctx context.Context, eventsArr []events.Event) {

	for _, event := range eventsArr {

		go func(ctx context.Context, event events.Event) {
			log := c.log.With(
				slog.String("event", event.Text),
			)

			tw := time.Now()

			log.Info("got new event")

			if err := c.processor.Process(ctx, event); err != nil {
				log.Error("can't handle event: ", sl.Err(err))
			}

			log.Debug("event is over", slog.Any("processing time", time.Since(tw)))

		}(ctx, event)
	}
}
