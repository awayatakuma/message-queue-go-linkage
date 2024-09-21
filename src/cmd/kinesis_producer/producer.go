package main

import (
	"context"
	"encoding/json"
	"fmt"
	"main/consts"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

func main() {

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(consts.REGION_NAME),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(_, _ string, _ ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:       consts.PARTITION_KEY,
				URL:               consts.ENDPOINT_URL_KINESIS,
				SigningRegion:     consts.REGION_NAME,
				HostnameImmutable: true,
			}, nil
		})),
	)
	if err != nil {
		panic(err)
	}

	kinesisClient := kinesis.NewFromConfig(cfg)
	if err != nil {
		panic(err)
	}

	fmt.Println("Programming is running... You can stop the process with Ctrl+C")
	i := 0
	loc, _ := time.LoadLocation("UTC")
	go func() {
		for {
			i += 1
			data := &consts.Data{
				HostName:  "kinesis-publisher",
				TimeStamp: time.Now().Local().In(loc).Format(time.RFC3339Nano),
				Id:        i,
				Payload:   "Message " + strconv.Itoa(i),
			}
			json, err := json.Marshal(data)
			if err != nil {
				fmt.Printf("Error: %s", err)
				return
			}
			fmt.Printf("inbound val: %v\n", string(json))
			var output *kinesis.PutRecordOutput
			pri := kinesis.PutRecordInput{
				Data:         json,
				StreamName:   aws.String(consts.STREAM_NAME),
				PartitionKey: aws.String(consts.PARTITION_KEY),
			}
			pri.StreamARN = aws.String(consts.STREAM_ARN)
			output, err = kinesisClient.PutRecord(
				ctx,
				&pri,
			)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%v\n", output)
		}
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	fmt.Println("Ctrl+C is pressed... The program is shutting down.")

}
