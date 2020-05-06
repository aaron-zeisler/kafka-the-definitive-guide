package main

import (
	"fmt"
	"math/rand"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/urfave/cli"
)

func customPartitioner(c *cli.Context) {
	topicName := "custom-partitioner-test"
	possibleKeys := []string{"important", "normal"}

	// Create a producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer producer.Close()

	// Create a channel to receive delivery reports
	deliveryReportsChannel := make(chan kafka.Event)

	for i := 0; i < 10; i++ {
		// Randomly choose one of the possible keys
		key := []byte(possibleKeys[rand.Int31n(int32(len(possibleKeys)))])

		// Choose the partition based on the key
		partition, err := choosePartition(producer, topicName, key)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Chose partition %d for the %s message\n", partition, string(key))

		// Define the message
		message := kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topicName,
				Partition: partition,
			},
			Key:   key,
			Value: []byte("test"),
		}

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
			return
		} else {
			fmt.Printf("Successfully produced a message and verified that it was sent: %#v\n", m)
		}
	}
}

func choosePartition(producer *kafka.Producer, topicName string, key []byte) (int32, error) {
	metadata, err := producer.GetMetadata(&topicName, true, 100)
	if err != nil {
		return 0, fmt.Errorf("Failed to get metadata from the broker: %w", err)
	}

	numPartitions := int32(len(metadata.Topics[topicName].Partitions))

	if string(key) == "important" {
		return numPartitions - 1, nil
	}
	return rand.Int31n(numPartitions - 1), nil
}
