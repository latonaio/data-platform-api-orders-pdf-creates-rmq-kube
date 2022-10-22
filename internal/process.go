package internal

import (
	configs "data-platform-api-orders-pdf-creates-rmq-kube/configs"
	"data-platform-api-orders-pdf-creates-rmq-kube/formatter"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/latonaio/golang-logging-library-for-sap/logger"
	rabbitmqClient "github.com/latonaio/rabbitmq-golang-client"
)

var mountPdfPath string
var rmqQueueTo string
var l *logger.Logger

const (
	CREATES = "creates"
	READS   = "reads"
	UPDATES = "updates"
	DELETES = "deletes"
)

func init() {
	cfgs, err := configs.New()
	if err != nil {
		panic(err)
	}
	mountPdfPath = cfgs.MountPdfPath
	rmqQueueTo = cfgs.RMQ.QueueTo()[0]

	l = logger.NewLogger()
}

func CallProcess(rmq *rabbitmqClient.RabbitmqClient, msg rabbitmqClient.RabbitmqMessage) (err error) {
	raw, err := json.Marshal(msg.Data())
	if err != nil {
		l.Fatal("data marshal error", err)
		return err
	}

	inputServiceParam := formatter.InputServiceParam{}
	err = json.Unmarshal(raw, &inputServiceParam)
	if err != nil {
		l.Fatal("inputServiceParam unmarshal error", err)
		return err
	}

	// it must return businessPartner orderId pdfData
	rmqSendParam, filename, pdfData := createResponse(
		raw,
		inputServiceParam,
	)

	l.Info("runtime session id: ", inputServiceParam.RuntimeSessionId)

	if err != nil {
		l.Fatal("data unmarshal error", err)
		return err
	}

	_, err = generatePdf(filename, pdfData)
	if err != nil {
		l.Fatal("failed to generate pdf", err)
		return err
	}

	err = rmq.Send(rmqQueueTo, rmqSendParam)

	l.Info("complete to sending queue: ", rmqQueueTo)

	if err != nil {
		l.Fatal("failed to send queue message", err)
		return err
	}
	return nil
}

func createResponse(raw []byte, inputServiceParam formatter.InputServiceParam) (map[string]interface{}, string, string) {
	apiServiceName := inputServiceParam.APIServiceName
	apiType := inputServiceParam.APIType

	if apiServiceName == formatter.DPFM_API_ORDERS_SRV && apiType == CREATES {
		return formatter.CreateResponseOrdersHeaderPDF(raw, inputServiceParam)
	}

	return nil, "", ""
}

func generatePdf(filename string, pdfData string) (string, error) {
	dec, err := base64.StdEncoding.DecodeString(pdfData)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf(filename)
	filePath := fmt.Sprintf("%s/%s",
		mountPdfPath,
		fileName)
	if err != nil {
		return "", err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return "", err
	}
	if err := f.Sync(); err != nil {
		return "", err
	}

	return fileName, nil
}
