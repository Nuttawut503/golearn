package main

import (
	"context"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// for initializing topics
func RunAdmin(servers string) {
	admin, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": servers,
	})
	if err != nil {
		log.Fatal("Admin can't connect")
	}
	defer admin.Close()

	var topicSpecifications []kafka.TopicSpecification

	for _, topic := range Topics {
		topicSpecifications = append(topicSpecifications, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     2,
			ReplicationFactor: 1,
		})
	}

	results, err := admin.CreateTopics(context.Background(), topicSpecifications, kafka.SetAdminOperationTimeout(time.Second*5))

	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		log.Printf("%s\n", result)
	}
}
