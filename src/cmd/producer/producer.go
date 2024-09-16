package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/consts"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

func main() {
	consts.Setup()
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           consts.ENDPOINT_URL_SNS,
					SigningRegion: "ap-northeast-1",
				}, nil
			}),
		),
	)
	if err != nil {
		log.Fatalf("LoadDefaultConfig failed:%v", err)
	}
	// Create SNS client
	svc := sns.NewFromConfig(cfg)

	fmt.Println("プログラムが開始されました。Ctrl+Cを押すと終了します。")
	i := 0
	loc, _ := time.LoadLocation("UTC")
	go func() {
		for {
			i++
			time.Sleep(100_000_000)
			// Publish a message to a topic
			data := &consts.Data{
				HostName:  "go-sns-sqs-go-publisher",
				TimeStamp: time.Now().Local().In(loc).Format(time.RFC3339Nano),
				Id:        i,
				Payload:   "Message " + strconv.Itoa(i),
			}
			msg, err := json.Marshal(data)
			if err != nil {
				log.Fatalf("failed to marshal: %v", err)
			}
			fmt.Printf("inbound val: %v\n", string(msg))
			result, err := svc.Publish(
				ctx,
				&sns.PublishInput{
					Message:        aws.String(string(msg)),
					TopicArn:       aws.String(consts.TOPIC_ARN),
					MessageGroupId: aws.String("singleton"),
				})
			if err != nil {
				log.Fatalf("Error publishing message: %v", err)
			}

			log.Println("Message published:", *result.MessageId)

		}
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	fmt.Println("Ctrl+Cが押されました。プログラムを終了します。")
}
