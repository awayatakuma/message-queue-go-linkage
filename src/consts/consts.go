package consts

import (
	"fmt"
	"os"
)

var (
	ENV = os.Getenv("ENV")

	ENDPOINT_URL_SNS        string
	ENDPOINT_URL_SQS        string
	ENDPOINT_URL_LOCAL      = "http://localstack:4566"
	ENDPOINT_URL_SNS_REMOTE = "https://sns.ap-northeast-1.amazonaws.com"
	ENDPOINT_URL_SQS_REMOTE = "https://sqs.ap-northeast-1.amazonaws.com"
	// SNS
	TOPIC_ARN        string
	TOPIC_ARN_LOCAL  = "arn:aws:sns:ap-northeast-1:000000000000:dummy-topic.fifo"
	TOPIC_ARN_REMOTE = "arn:aws:sns:ap-northeast-1:553695339919:MySnsSqsFanoutTopic.fifo"

	// SQS
	QUEUE_URL        string
	QUEUE_URL_LOCAL  = "http://localstack:4566/000000000000/dummy-queue.fifo"
	QUEUE_URL_REMOTE = "https://sqs.ap-northeast-1.amazonaws.com/553695339919/MySqsQueue-0.fifo"

	LOG_DIR = "./logs/"
)

func Setup() {
	fmt.Printf("Environment: %s\n", ENV)
	switch ENV {
	case "local":
		ENDPOINT_URL_SNS = ENDPOINT_URL_LOCAL
		ENDPOINT_URL_SQS = ENDPOINT_URL_LOCAL

		TOPIC_ARN = TOPIC_ARN_LOCAL
		QUEUE_URL = QUEUE_URL_LOCAL
	case "remote":
		ENDPOINT_URL_SNS = ENDPOINT_URL_SNS_REMOTE
		ENDPOINT_URL_SQS = ENDPOINT_URL_SQS_REMOTE

		TOPIC_ARN = TOPIC_ARN_REMOTE
		QUEUE_URL = QUEUE_URL_REMOTE
	default:
		panic("invalid environment name. please set either local or remote in .env file")
	}
}

type Data struct {
	HostName  string
	TimeStamp string
	Id        int
	Payload   string
}
