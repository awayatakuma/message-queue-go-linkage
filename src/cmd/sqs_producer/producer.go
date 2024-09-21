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
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           consts.QUEUE_URL,
					SigningRegion: "ap-northeast-1",
				}, nil
			}),
		),
	)
	if err != nil {
		log.Fatalf("LoadDefaultConfig failed:%v", err)
	}
	// Create SNS client
	svc := sqs.NewFromConfig(cfg)

	fmt.Println("Programming is running... You can stop the process with Ctrl+C")
	i := 0
	loc, _ := time.LoadLocation("UTC")
	go func() {
		for {
			i++
			time.Sleep(100_000_000)
			data := &consts.Data{
				HostName:  "sqs-publisher",
				TimeStamp: time.Now().Local().In(loc).Format(time.RFC3339Nano),
				Id:        i,
				Payload:   "Message " + strconv.Itoa(i),
			}
			msg, err := json.Marshal(data)
			if err != nil {
				log.Fatalf("failed to marshal: %v", err)
			}
			fmt.Printf("inbound val: %v\n", string(msg))
			result, err := svc.SendMessage(
				ctx,
				&sqs.SendMessageInput{
					MessageBody:    aws.String(string(msg)),
					QueueUrl:       aws.String(consts.QUEUE_URL),
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
	fmt.Println("Ctrl+C is pressed... The program is shutting down.")
}
