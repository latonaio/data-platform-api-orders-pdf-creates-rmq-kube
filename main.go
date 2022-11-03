package main

import (
	"data-platform-api-orders-pdf-creates-rmq-kube/configs"
	"data-platform-api-orders-pdf-creates-rmq-kube/internal"

	"github.com/latonaio/golang-logging-library-for-data-platform/logger"
	rabbitmq "github.com/latonaio/rabbitmq-golang-client-for-data-platform"
)

func main() {
	l := logger.NewLogger()
	cfgs, err := configs.New()
	if err != nil {
		l.Fatal("failed to get kanban client", err)
	}

	l.Info("created configs instance", cfgs.RMQ.URL(), cfgs.RMQ.QueueFrom(), cfgs.RMQ.QueueTo())

	rmq, err := rabbitmq.NewRabbitmqClient(cfgs.RMQ.URL(), cfgs.RMQ.QueueFrom(), "", cfgs.RMQ.QueueTo(), -1)

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
