package main

import (
	"fmt"

	"github.com/tech-scripter/sandbox/logging/rlogging"
	"github.com/tech-scripter/sandbox/logging/slogging"
)

var (
	err = fmt.Errorf("error")
)

func main() {
	s := slogging.New()
	r := rlogging.New()

	r.
		With("amqp", "queue", "deadLetterQueueName", "exchange", "deadLetterExchangeName").
		Debug("AMQP message logged")

	fmt.Println()

	r.With("", "key", "value").Debug("testing")

	fmt.Println()

	s.
		With("amqp", "queue", "deadLetterQueueName", "exchange", "deadLetterExchangeName").
		Debug("AMQP message logged")

	fmt.Println()

	s.
		With("", "key", "value").Debug("testing")
}
