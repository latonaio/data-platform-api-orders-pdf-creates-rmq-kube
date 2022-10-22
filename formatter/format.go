package formatter

import (
	"github.com/latonaio/golang-logging-library-for-sap/logger"
)

type InputServiceParam struct {
	APIServiceName   string `json:"aPIServiceName"`
	APIType          string `json:"aPIType"`
	RuntimeSessionId string `json:"runtime_session_id"`
}

type InputParam struct {
	InputOrdersHeaderPDF
}

var l *logger.Logger

func init() {
	l = logger.NewLogger()
}
