package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/config"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
	"github.com/markusylisiurunen/template-monorepo/package/go/hello"
)

func sayHelloRepeatedly(ctx context.Context, done chan bool, name string) {
	ticker := time.NewTicker(5 * time.Second)

	logger.Default.Infof("starting to say hello to \"%s\"...", name)

	for {
		select {
		case <-ticker.C:
			hello.Say(name, func(msg string) { logger.Default.Infof(msg) })
		case <-ctx.Done():
			logger.Default.Infof("finished saying hello!")
			done <- true

			return
		}
	}
}

func main() {
	logger.Default.Infof("starting the service...")

	if err := config.Load(); err != nil {
		logger.Default.Errorf(err.Error())
		os.Exit(1)
	}

	cfg, err := config.Get()
	if err != nil {
		os.Exit(1)
	}

	logger.Default.Infof("config loaded successfully")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan bool, 1)

	go sayHelloRepeatedly(ctx, done, "Swiftbeaver")

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT)
	signal.Notify(signals, syscall.SIGTERM)

	<-signals

	logger.Default.Infof("shutting down...")

	cancel()
	<-done

	logger.Default.Infof("all done!")
}
