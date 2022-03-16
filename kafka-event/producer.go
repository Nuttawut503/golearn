package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ProduceMessage(producer *kafka.Producer, topic string, value []byte) error {
	deliveryChan := make(chan kafka.Event)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}, deliveryChan)
	if err != nil {
		return err
	}
	event := <-deliveryChan
	if message, ok := event.(*kafka.Message); ok {
		if message.TopicPartition.Error != nil {
			log.Printf("Produce failed: %v\n", message.TopicPartition.Error)
			return message.TopicPartition.Error
		} else {
			log.Printf(
				"Produce succeeded: topic - %v + partition - %v + offset - %v\n",
				*message.TopicPartition.Topic, message.TopicPartition.Partition, message.TopicPartition.Offset,
			)
			return nil
		}
	}
	return errors.New("unknown error")
}

func RegisterAccountHandler(producer *kafka.Producer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var newAccount RegisterAccountEvent
		if err := ctx.ShouldBindJSON(&newAccount); err != nil {
			ctx.String(http.StatusUnprocessableEntity, "Wrong format")
			return
		}
		newAccount.TransactionID = uuid.New().String()
		topic := reflect.TypeOf(newAccount).Name()
		value, _ := json.Marshal(newAccount)
		if err := ProduceMessage(producer, topic, value); err != nil {
			ctx.String(http.StatusInternalServerError, "Something went wrong")
			return
		}
		ctx.String(http.StatusOK, "Registration request has been processed")
	}
}

func DeactivateAccountHandler(producer *kafka.Producer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var account DeactivateAccountEvent
		if err := ctx.ShouldBindJSON(&account); err != nil {
			ctx.String(http.StatusUnprocessableEntity, "Wrong format")
			return
		}
		account.TransactionID = uuid.New().String()
		topic := reflect.TypeOf(account).Name()
		value, _ := json.Marshal(account)
		if err := ProduceMessage(producer, topic, value); err != nil {
			ctx.String(http.StatusInternalServerError, "Something went wrong")
			return
		}
		ctx.String(http.StatusOK, "Deactivation request has been processed")
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	r := gin.Default()

	r.POST("/account", RegisterAccountHandler(producer))
	r.DELETE("/account", DeactivateAccountHandler(producer))

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()

	stop()
	log.Println("shutting down gracefully")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
