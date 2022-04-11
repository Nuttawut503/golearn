package main

import (
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func ProduceMessage(producer *kafka.Producer, topic string, value []byte) error {
	return producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, nil)
}

func eventListener(events chan kafka.Event) {
	for ev := range events {
		m, ok := ev.(*kafka.Message)
		if !ok {
			continue
		}
		if m.TopicPartition.Error != nil {
			fmt.Fprintf(os.Stderr, "%% Delivery error: %v\n", m.TopicPartition)
		} else {
			fmt.Fprintf(os.Stderr, "%% Delivered %v\n", m)
		}
	}
}

func RunProducer(servers string) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": servers,
	})
	if err != nil {
		log.Fatal("Producer can't connect")
	}

	defer producer.Close()

	go eventListener(producer.Events())

	runServerWithGracefullyShutdown(producer)
}
