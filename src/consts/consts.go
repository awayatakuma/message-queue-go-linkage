package consts

var (
	ENDPOINT_URL_SQS = "https://sqs.ap-northeast-1.amazonaws.com"

	// SQS
	QUEUE_URL = "https://sqs.ap-northeast-1.amazonaws.com/826758858441/MyQueue.fifo"

	LOG_DIR = "./logs/"

	ENDPOINT_URL_KINESIS  = "https://kinesis.ap-northeast-1.amazonaws.com"
	ENDPOINT_URL_DYNAMODB = "https://dynamodb.ap-northeast-1.amazonaws.com"

	// Kinesis
	STREAM_ARN    = "arn:aws:kinesis:ap-northeast-1:826758858441:stream/MyKinesis"
	SHARD_ID      = "shardId-000000000000"
	PARTITION_KEY = "my/topic"
	REGION_NAME   = "ap-northeast-1"
	STREAM_NAME   = "MyKinesis"
)

type Data struct {
	HostName  string
	TimeStamp string
	Id        int
	Payload   string
}
