package main

import (
	"fmt"

	"github.com/tech-scripter/sandbox/logging"
	"github.com/tech-scripter/sandbox/logging/rlogging"
	"github.com/tech-scripter/sandbox/logging/slogging"
)

func main() {
	rlog := rlogging.New()
	slog := slogging.New()

	doSmth(rlog)
	doSmth(slog)
}

func doSmth(log logging.Logger) {
	var err = fmt.Errorf("test error %s", "test")
	log.WithError(err).Error("test error")
}
