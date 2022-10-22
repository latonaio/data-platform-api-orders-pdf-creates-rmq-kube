package main

import (
	"data-platform-api-orders-pdf-creates-rmq-kube/configs"
	"data-platform-api-orders-pdf-creates-rmq-kube/internal"

	"github.com/latonaio/golang-logging-library-for-sap/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client"
)

func main() {
	l := logger.NewLogger()
	cfgs, err := configs.New()
	if err != nil {
		l.Fatal("failed to get kanban client", err)
	}

	l.Info("created configs instance", cfgs.RMQ.URL(), cfgs.RMQ.QueueFrom(), cfgs.RMQ.QueueTo())

	rmq, err := rabbitmq.NewRabbitmqClient(cfgs.RMQ.URL(), cfgs.RMQ.QueueFrom(), cfgs.RMQ.QueueTo())

	if err != nil {
		l.Fatal(err.Error())
	}
	defer rmq.Close()

	iter, err := rmq.Iterator()
	if err != nil {
		l.Fatal(err.Error())
	}
	defer rmq.Stop()

	for msg := range iter {
		l.Info("received queue message")
		err = internal.CallProcess(rmq, msg)
		if err != nil {
			msg.Fail()
			l.Info("call process error", err)
			continue
		}
		msg.Success()
	}
}
