# message-queue-go-linkage

This repository is created to confirm the performance of SQS and Kinesis Data Stream of AWS.  
I summed up the result in [this blog article](https://asynctp.tokyo/blog/what-you-really-need-may-sqs) so please read it if you feel like to know it.

## Prerequirements
Please prepare AWS environments by yourself.  
I prepared Iac(CDK) for this verifications: https://github.com/awayatakuma/cdk-sqs-kinesis-performance-check  
It is easiest way to create the environment to use the above codes.  
Please read [README.md](https://github.com/awayatakuma/cdk-sqs-kinesis-performance-check/blob/main/README.md) if you feel like to know how to use the codes.

## How to use
- Access to EC2 instance by ssh
- Do `docker compose up -d`
- Open two terminals

### SQS
- Consumers
  - Activate the process by: 
    - `docker exec -it go-awsclient go run cmd/sqs_consumer/consumer.go`
- Producers
  - After confirming consumers processes activated, activate producers processes by: 
    - `docker exec -it go-awsclient go run cmd/sqs_producer/producer.go`
- You can check pub-sub processes in your terminal. If you feel enough data are gathered, deactivate the processes by Ctrl-C`
- The results can be confirmed in `logs` folders

### Kinesis
- Consumers
  - Activate the process by: 
    - `docker exec -it go-awsclient go run cmd/kinesis_consumer/consumer.go`
- Producers
  - After confirming consumers processes activated, activate producers processes by: 
    - `docker exec -it go-awsclient go run cmd/kinesis_producer/producer.go`
- You can check pub-sub processes in your terminal. If you feel enough data are gathered, deactivate the processes by Ctrl-C
- The results can be confirmed in `logs` folders