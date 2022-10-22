package formatter

import (
	"encoding/json"
)

const DPFM_API_ORDERS_SRV = "DPFM_API_ORDERS_SRV"
const function = "OrdersHeaderPDF"

type InputOrdersHeaderPDF struct {
	ConnectionKey string `json:"connection_key"`
	Result        bool   `json:"result"`
	RedisKey      string `json:"redis_key"`
	Filepath      string `json:"filepath"`
	Orders        struct {
		OrderId         int `json:"OrderID"`
		BusinessPartner int `json:"BusinessPartner"`
		HeaderPDF       struct {
			DocType      string `json:"DocType"`
			DocVersionID int    `json:"DocVersionID"`
			DocID        string `json:"DocID"`
			FileName     string `json:"FileName"`
		}
	}
	APISchema     string   `json:"api_schema"`
	Accepter      []string `json:"accepter"`
	PDFData       string   `json:"pdf_data"`
	ValidatedDate string   `json:"validated_date"`
	Deleted       bool     `json:"deleted"`
	OrderId       string   `json:"older_id"`
}

type ResponseOrdersHeaderPDF struct {
	BusinessPartner int    `json:"BusinessPartner"`
	OrderID         int    `json:"OrderID"`
	DocType         string `json:"DocType"`
	DocVersionID    int    `json:"DocVersionID"`
	DocID           string `json:"DocID"`
	FileName        string `json:"FileName"`
}

func CreateResponseOrdersHeaderPDF(raw []byte, inputServiceParam InputServiceParam) (map[string]interface{}, string, string) {
	response := ResponseOrdersHeaderPDF{}
	err := json.Unmarshal(raw, &response)
	if err != nil {
		l.Fatal("response unmarshal error", err)
	}

	orderParam := InputOrdersHeaderPDF{}
	err = json.Unmarshal(raw, &orderParam)

	var responses []ResponseOrdersHeaderPDF
	responses = append(responses, ResponseOrdersHeaderPDF{
		BusinessPartner: orderParam.Orders.BusinessPartner,
		OrderID:         orderParam.Orders.OrderId,
		DocType:         orderParam.Orders.HeaderPDF.DocType,
		DocVersionID:    orderParam.Orders.HeaderPDF.DocVersionID,
		DocID:           orderParam.Orders.HeaderPDF.DocID,
		FileName:        orderParam.Orders.HeaderPDF.FileName,
	})

	rmqQueueFormat := map[string]interface{}{
		"runtime_session_id": inputServiceParam.RuntimeSessionId,
		"message":            responses,
		"function":           function,
	}

	l.Info(rmqQueueFormat)

	return rmqQueueFormat,
		orderParam.Orders.HeaderPDF.FileName,
		orderParam.PDFData
}
