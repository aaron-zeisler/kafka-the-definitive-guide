package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/urfave/cli"
)

func asynchronous(c *cli.Context) {
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

	// Produce the message
	err = producer.Produce(&message, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("The message was sent asynchronously!")

	// Listen on the producer's delivery reports channel ...
	e := <-producer.Events()

	// The thing that came through the delivery reports channel should be a kafka.Message
	m := e.(*kafka.Message)
	// Check message.TopicPartition.Error to see if an error happened
	if m.TopicPartition.Error != nil {
		fmt.Printf("The message failed to be produced: %s\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Verified that the message was sent: %#v\n", m)
	}
}
