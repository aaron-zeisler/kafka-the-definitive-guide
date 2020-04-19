package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/urfave/cli"
)

func synchronous(c *cli.Context) {
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
		Value: []byte("synchronous with the Confluent client"),
	}

	// Create a channel to receive delivery reports
	deliveryReportsChannel := make(chan kafka.Event)

	// Produce the message
	err = producer.Produce(&message, deliveryReportsChannel)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Listen on the delivery reports channel ...
	e := <-deliveryReportsChannel

	// The thing that came through the delivery reports channel should be a kafka.Message
	m := e.(*kafka.Message)
	// Check message.TopicPartition.Error to see if an error happened
	if m.TopicPartition.Error != nil {
		fmt.Printf("The message failed to be produced: %s\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Successfully produced a message and verified that it was sent: %#v\n", m)
	}
}
