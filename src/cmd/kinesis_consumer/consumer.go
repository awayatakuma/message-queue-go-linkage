package main

import (
	"encoding/json"
	"fmt"
	"log"
	"main/consts"
	"main/time_stamp_logger"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vmware/vmware-go-kcl-v2/clientlibrary/config"
	"github.com/vmware/vmware-go-kcl-v2/clientlibrary/interfaces"
	"github.com/vmware/vmware-go-kcl-v2/clientlibrary/worker"
)

type RecordProcessorFactory struct{}

func (r *RecordProcessorFactory) CreateProcessor() interfaces.IRecordProcessor {
	return &RecordProcessor{}
}

type RecordProcessor struct {
}

func (r *RecordProcessor) Initialize(input *interfaces.InitializationInput) {
	log.Printf("Initializing record processor for shard")
}

func (r *RecordProcessor) ProcessRecords(input *interfaces.ProcessRecordsInput) {
	log.Printf("Processing %d record(s)", len(input.Records))
	for _, record := range input.Records {

		var data consts.Data
		if err := json.Unmarshal(record.Data, &data); err != nil {
			panic(err)
		}
		fmt.Printf("Received data :%+v\n ", data)
		start, _ := time.Parse(time.RFC3339Nano, data.TimeStamp)
		time_stamp_logger.Write(data.Id, start, time.Now())
	}
}

func (r *RecordProcessor) Shutdown(input *interfaces.ShutdownInput) {
	log.Printf("Shutting down record processor for shard")
}

func main() {

	applicationName := "DynamoTable" // dynamoで使う
	workerID := "my-worker"          // kcl内部でのみ使う

	time_stamp_logger.Initial("kcl")

	kclConfig := config.NewKinesisClientLibConfig(
		applicationName,
		consts.STREAM_NAME,
		consts.REGION_NAME,
		workerID,
	).
		WithInitialPositionInStream(config.LATEST).
		WithTableName(applicationName).
		WithMaxRecords(1).
		WithIdleTimeBetweenReadsInMillis(1).
		WithEnhancedFanOutConsumer(true).
		WithKinesisEndpoint(consts.ENDPOINT_URL_KINESIS).
		WithDynamoDBEndpoint(consts.ENDPOINT_URL_DYNAMODB)

	worker := worker.NewWorker(&RecordProcessorFactory{}, kclConfig)

	err := worker.Start()
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	fmt.Println("Programming is running... You can stop the process with Ctrl+C")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	fmt.Println("Ctrl+C is pressed... The program is shutting down.")
}
