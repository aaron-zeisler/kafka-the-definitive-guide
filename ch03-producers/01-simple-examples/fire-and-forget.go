package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/urfave/cli"
)

// fireAndForget method using the Confluent Kafka client
func fireAndForget(c *cli.Context) {
	// Create the producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer producer.Close()

	// Define the message
	topicName := "test"
	message := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &topicName,
		},
		Value: []byte("fire-and-forget with the Confluent client"),
	}

	// Produce the message
	err = producer.Produce(&message, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("successfully produced a message")
}
