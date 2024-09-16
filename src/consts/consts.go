package consts

var (
	ENDPOINT_URL_SQS = "https://sqs.ap-northeast-1.amazonaws.com"

	// SQS
	QUEUE_URL = "https://sqs.ap-northeast-1.amazonaws.com/826758858441/MyQueue.fifo"

	LOG_DIR = "./logs/"
)

type Data struct {
	HostName  string
	TimeStamp string
	Id        int
	Payload   string
}
