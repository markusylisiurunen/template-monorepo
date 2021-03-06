package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	opinionatedevents "github.com/markusylisiurunen/go-opinionated-events"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/config"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/entities/messages"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/events"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/logger"
	"github.com/markusylisiurunen/template-monorepo/backend/go-example/pkg/migrations"
	"github.com/markusylisiurunen/template-monorepo/package/go/hello"
)

func sayHelloRepeatedly(ctx context.Context, done chan bool, name string) {
	ticker := time.NewTicker(5 * time.Second)

	logger.Default.Infof("starting to say hello to \"%s\"...", name)

	for {
		select {
		case <-ticker.C:
			publisher, err := events.GetPublisher()
			if err != nil {
				logger.Default.Errorw(err.Error())
				continue
			}

			hello.Say(name, func(msg string) {
				logger.Default.Infow(msg)
			})

			message := opinionatedevents.NewMessage("hello")
			payload := messages.NewHelloPayload("Howdy", "Stranger")

			if err := message.SetPayload(payload); err != nil {
				logger.Default.Errorw(err.Error())
				continue
			}

			if err := publisher.Publish(message); err != nil {
				logger.Default.Errorw(err.Error())
				continue
			}
		case <-ctx.Done():
			logger.Default.Infow("finished saying hello!")
			done <- true

			return
		}
	}
}

func drainPublisher(ctx context.Context, done chan bool) {
	<-ctx.Done()

	publisher, err := events.GetPublisher()
	if err == nil {
		publisher.Drain()
	}

	done <- true
}

func startServer(ctx context.Context, done chan bool, cfg *config.Config) {
	router, err := setupHttpEndpoints(cfg, logger.Default)
	if err != nil {
		logger.Default.Errorf("failed to create endpoints")
		os.Exit(1)
	}

	// serve the documentation from the `static` directory
	router.
		PathPrefix("/_docs").
		Handler(
			http.StripPrefix(
				"/_docs/",
				http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
					resp.Header().Del("Content-Type")
					http.FileServer(http.Dir("./static")).ServeHTTP(resp, req)
				}),
			),
		)

	address := cfg.ServerHost + ":" + strconv.Itoa(cfg.ServerPort)
	srv := &http.Server{Addr: address}

	http.Handle("/", router)

	logger.Default.Infow("starting the server...",
		"Host", cfg.ServerHost,
		"Port", cfg.ServerPort,
	)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			logger.Default.Errorw("http server returned an error",
				"Error", err.Error(),
			)
		}
	}()

	<-ctx.Done()

	logger.Default.Infow("shutting down the server...")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Default.Errorw("http server could not be shut down gracefully",
			"Error", err.Error(),
		)
	}

	logger.Default.Infow("server shut down")

	done <- true
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

	if err := migrations.Migrate("migrations", cfg.DatabaseURL); err != nil {
		logger.Default.Errorw("could not run migrations",
			"Error", err.Error(),
		)

		os.Exit(1)
	}

	logger.Default.Infof("migrations run successfully")

	ctx, cancel := context.WithCancel(context.Background())

	sayHelloIntervalDone := make(chan bool, 1)
	drainPublisherDone := make(chan bool, 1)
	httpServerDone := make(chan bool, 1)

	go sayHelloRepeatedly(ctx, sayHelloIntervalDone, "Swiftbeaver")
	go drainPublisher(ctx, drainPublisherDone)
	go startServer(ctx, httpServerDone, cfg)

	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	<-signals

	logger.Default.Infof("shutting down...")

	cancel()

	<-sayHelloIntervalDone
	<-drainPublisherDone
	<-httpServerDone

	logger.Default.Infof("all done!")
}
