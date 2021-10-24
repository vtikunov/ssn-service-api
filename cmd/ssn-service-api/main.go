package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	retranslatorpkg "github.com/ozonmp/ssn-service-api/internal/app/retranslator"
)

func main() {
	sigs := make(chan os.Signal, 1)

	retranslator := retranslatorpkg.NewRetranslator(
		&retranslatorpkg.Configuration{
			EventChannelSize: 100,

			MaxConsumers:      10,
			ConsumerTimeout:   time.Millisecond * 500,
			ConsumerBatchTime: time.Millisecond * 100,
			ConsumerBatchSize: 500,

			MaxProducers:       10,
			ProducerTimeout:    time.Millisecond * 500,
			ProducerMaxWorkers: 2,
		},
	)
	retranslator.Start()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
