package main

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/urfave/cli"
)

func consumeSpecificPartition(c *cli.Context) {
	broker := "localhost:9092"
	// This example uses the topic that was created in chapter 3: custom-partitioner.
	//  In that example, 'important' messages were written to partition 4 and 'normal' messages were written to partitions 0-3.
	//  In this example, the consumer specified partition 4 so it can just read the 'important' messages.
	topicName := "custom-partitioner-test"
	var partitionNumber int32 = 4

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{broker}, config)
	if err != nil {
		fmt.Printf("failed to create the kafka consumer: %s", err)
		return
	}
	defer consumer.Close()

	// Create a client that consumes a single partition
	client, err := consumer.ConsumePartition(topicName, partitionNumber, 0)
	if err != nil {
		fmt.Printf("failed to connect to the partition: %s", err)
		return
	}

	// Listen for errors
	go func() {
		for {
			err := <-client.Errors()
			fmt.Printf("An error came through the Errors channel: %s\n", err)
		}
	}()

	// Consume messages from this partition
	for msg := range client.Messages() {
		fmt.Printf("received an 'important' message! Topic: %s; Partition: %d; Offset: %d; Key: %s; Value: %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}
}
