package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/mtaylor91/yakd/pkg/build/cmd"
	"github.com/mtaylor91/yakd/pkg/util/log"
)

func main() {
	ctx, log := log.Setup(context.Background())
	ctx, cancel := context.WithCancel(ctx)
	exitSignals := make(chan os.Signal, 1)
	signal.Notify(exitSignals, os.Interrupt, os.Kill)
	go waitForSignal(ctx, log, cancel, exitSignals)
	cmd.Root.ExecuteContext(ctx)
}

func waitForSignal(
	ctx context.Context,
	log *logrus.Entry,
	cancel context.CancelFunc,
	exitSignals chan os.Signal,
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
