package main

import (
	"context"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/cmd"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	exitSignals := make(chan os.Signal, 1)
	signal.Notify(exitSignals, os.Interrupt, os.Kill)
	go waitForSignal(ctx, cancel, exitSignals)
	cmd.Root.ExecuteContext(ctx)
}

func waitForSignal(
	ctx context.Context, cancel context.CancelFunc, exitSignals chan os.Signal,
) {
	select {
	case <-exitSignals:
		log.Warnf("Received exit signal")
		cancel()
	case <-ctx.Done():
		break
	}

	// Force exit on second signal
	<-exitSignals
	log.Errorf("Received second exit signal, forcing exit")

	os.Exit(1)
}
