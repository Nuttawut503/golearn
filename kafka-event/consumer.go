package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func handleEvent(topic *string, value []byte) {
	if topic == nil {
		return
	}
	log.Println(*topic)
	switch *topic {
	case reflect.TypeOf(RegisterAccountEvent{}).Name():
		var newAccount RegisterAccountEvent
		if err := json.Unmarshal(value, &newAccount); err != nil {
			return
		}
		// do something
		log.Printf("Welcome %s (%s %s)\n", newAccount.Email, newAccount.Firstname, newAccount.Lastname)
	case reflect.TypeOf(DeactivateAccountEvent{}).Name():
		var goneAccount DeactivateAccountEvent
		if err := json.Unmarshal(value, &goneAccount); err != nil {
			return
		}
		// do something
		log.Printf("Goodbye %s", goneAccount.Email)
	}
}

func RunConsumer(servers string) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               servers,
		"group.id":                        "mygroup",
		"auto.offset.reset":               "earliest",
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"session.timeout.ms":              6000,
	})
	if err != nil {
		log.Fatal("Consumer can't connect")
	}
	defer consumer.Close()
	consumer.SubscribeTopics(Topics, nil)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sig:
			log.Println("Shutting down consumer")
			return
		case ev := <-consumer.Events():
			switch ev := ev.(type) {
			case kafka.AssignedPartitions:
				log.Println(ev)
				consumer.Assign(ev.Partitions)
			case kafka.RevokedPartitions:
				log.Println(ev)
				consumer.Unassign()
			case *kafka.Message:
				log.Printf("Message received: partition - %v + value: %v\n", ev.TopicPartition, string(ev.Value))
				handleEvent(ev.TopicPartition.Topic, ev.Value)
			case kafka.Error:
				log.Printf("Error occurred: %v", ev)
			}
		}
	}
}
