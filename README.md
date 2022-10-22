# data-platform-api-orders-pdf-creates-rmq-kube
data-platform-api-orders-pdf-creates-rmq-kubeは、データ連携基盤において、オーダーのPDFデータを登録・生成するマイクロサービスです。

## 動作環境    
* OS: Linux OS  
* CPU: ARM/AMD/Intel  


## PDFの生成
order.go の以下の箇所が、PDFデータを生成するソースコードです。

```
func generatePdf(orderParam OrderParam) (string, error) {
	dec, err := base64.StdEncoding.DecodeString(orderParam.PDFData)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s-%s-%s.pdf",
		strconv.Itoa(orderParam.Orders.BusinessPartner),
		strconv.Itoa(orderParam.Orders.OrderId),
		time.Now().Format("20060102150405"))
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
```

## RabbitMQ からの JSON Input
data-platform-api-orders-pdf-creates-rmq-kube は、入力ファイルとして、RabbitMQ からのメッセージを JSON 形式で受け取ります。

## RabbitMQ からのメッセージ受信による イベントドリヴン の ランタイム実行
data-platform-api-orders-pdf-creates-rmq-kube は、RabbitMQ からのメッセージを受け取ると、イベントドリヴンでランタイムを実行します。  
AION の仕様では、Kubernetes 上 の 当該マイクロサービスPod は 立ち上がったまま待機状態で当該メッセージを受け取り、（コンテナ起動などの段取時間をカットして）即座にランタイムを実行します。　 

## RabbitMQ の マスタサーバ環境
data-platform-api-orders-pdf-creates-rmq-kube が利用する RabbitMQ のマスタサーバ環境は、rabbitmq-on-kubernetes です。  

## RabbitMQ の Golang Runtime ライブラリ
data-platform-api-orders-pdf-creates-rmq-kube は、RabbitMQ の Golang Runtime ライブラリ として、rabbitmq-golang-clientを利用しています。

## デプロイ・稼働
data-platform-api-orders-pdf-creates-rmq-kube の デプロイ・稼働 を行うためには、aion-service-definitions の services.yml に、本レポジトリの services.yml を設定する必要があります。

kubectl apply - f 等で Deployment作成後、以下のコマンドで Pod が正しく生成されていることを確認してください。

```
$ kubectl get pods
```

## Output

data-platform-api-orders-pdf-creates-rmq-kube では、[golang-logging-library](https://github.com/latonaio/golang-logging-library) により、Outputとして、RabbitMQからのメッセージをJSON形式で出力します。"cursor" ～ "time"は、golang-logging-library による 定型フォーマットの出力結果です。

```
{
  "cursor": "/go/src/github.com/latonaio/formatter/OrdersHeaderPDF.go#L68",
  "function": "data-platform-api-orders-creates-pdf/formatter.CreateResponseOrdersHeaderPDF",
  "level": "INFO",
  "message": {
    "function": "OrdersHeaderPDF",
    "message": [
      {
        "BusinessPartner": 101,
        "OrderID": 1,
        "DocType": "ORDER",
        "DocVersionID": 1,
        "DocID": "92q9o0fa90g\u0026awp090aqqm0qopgwaepw9ere0weg",
        "FileName": "ラトナ株式会社殿_注文書_20221020_0001.pdf"
      }
    ],
    "runtime_session_id": "cd6d58b1-163c-4245-b6d0-5e87fc627807"
  },
  "time": "2022-10-21T08:47:29Z"
}

```