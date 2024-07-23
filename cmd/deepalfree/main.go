package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/al-ship/deepalfree"
)

var zapLogger *zap.Logger

type Service interface {
	Start() error
	Stop()
}

func main() {
	var err error

	zapLogger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	serv := deepalfree.NewHttp("ashipulya.ddns.net", zapLogger)
	httpClosedUnexpectedly := make(chan struct{})
	go func() {
		err := serv.Start()
		if err != nil {
			zapLogger.Error("start http server failed", zap.Error(err))
		}
		close(httpClosedUnexpectedly)
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-signals:
		zapLogger.Info("Signal received", zap.String("signal", s.String()))
		serv.Stop()
	case <-httpClosedUnexpectedly:
	}
}
