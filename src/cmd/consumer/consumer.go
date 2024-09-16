package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/consts"
	"main/time_stamp_logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	consts.Setup()
	time_stamp_logger.Initial("sdk")
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           consts.ENDPOINT_URL_SQS,
					SigningRegion: "ap-northeast-1",
				}, nil
			}),
		),
	)
	if err != nil {
		log.Fatalf("LoadDefaultConfig failed:%v", err)
	}
	// Create SNS client
	sqc := sqs.NewFromConfig(cfg)

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		fmt.Printf("failed to set locale.")
		return
	}
	fmt.Println("プログラムが開始されました。Ctrl+Cを押すと終了します。")
	i := 0
	go func() {
		input := &sqs.ReceiveMessageInput{
			QueueUrl: aws.String(consts.QUEUE_URL),
			// チューニングポイント
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     20,
		}
		for {
			i += 1
			res, err := sqc.ReceiveMessage(ctx, input)
			if err != nil {
				fmt.Printf("Error: %s\n", err)
				return

			}
			if len(res.Messages) <= 0 {
				log.Println("No Message Contains")
				continue
			}

			for _, msg := range res.Messages {
				_, err := sqc.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(consts.QUEUE_URL),
					ReceiptHandle: msg.ReceiptHandle,
				})
				if err != nil {
					panic(fmt.Sprintf("fail to delete msg: %v\n", err))
				}
				jsonMap := make(map[string]interface{})
				err = json.Unmarshal([]byte(*msg.Body), &jsonMap)
				if err != nil {
					panic(err)
				}
				var data consts.Data

				if err := json.Unmarshal([]byte(fmt.Sprintf("%v", jsonMap["Message"])), &data); err != nil {
					panic(fmt.Sprintf("fail to unmarshal msg: %v\n", err))
				}
				fmt.Printf("Unmarshaled data:%+v\n", data)
				stop := time.Now().Local().In(loc)
				start, err := time.Parse(time.RFC3339Nano, data.TimeStamp)
				if err != nil {
					fmt.Printf("failed to parse string to time")
					break
				}
				time_stamp_logger.Write(data.Id, start, stop)

			}

		}
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	fmt.Println("Ctrl+Cが押されました。プログラムを終了します。")
	time_stamp_logger.Stop()
}
