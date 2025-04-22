package main

import (
	"fmt"

	"github.com/tech-scripter/sandbox/logging"
	"github.com/tech-scripter/sandbox/logging/rlogging"
	"github.com/tech-scripter/sandbox/logging/slogging"
)

func main() {
	// concrete implementations
	rlog := rlogging.New()
	slog := slogging.New()

	doSmth(rlog)
	doSmth(slog)
}

func doSmth(log logging.Logger) {
	var err = fmt.Errorf("test error %s", "test")

	// Since the variable log is of interface type,
	// only the methods defined in the interface are accessible.
	// Therefore, the interface must include all the methods
	// that are used across our services.
	log.WithError(err).Error("test error")
}
