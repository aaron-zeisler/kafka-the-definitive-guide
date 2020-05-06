package main

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/urfave/cli"
)

// This simpleHandler implements the sarama.ConsumerGroupHandler interface.
// A ConsumerGroupHandler must be passed to sarama's Consume() method.
type simpleHandler struct{}

func (h simpleHandler) Setup(session sarama.ConsumerGroupSession) error   { return nil }
func (h simpleHandler) Cleanup(session sarama.ConsumerGroupSession) error { return nil }
func (h simpleHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("Received a message! Topic: %s; Partition: %d; Offset: %d; Key: %s; Value: %s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
	}
	return nil
}

func simpleConsumer(c *cli.Context) {
	broker := "localhost:9092"
	topicName := "test"
	consumerGroupID := "simple-consumer"

	config := sarama.NewConfig()
	config.Version = sarama.V2_0_0_0
	config.Consumer.Return.Errors = true

	// Create a consumer and add it to a consumer group
	consumer, err := sarama.NewConsumerGroup([]string{broker}, consumerGroupID, config)
	if err != nil {
		fmt.Printf("Failed to create the kafka consumer: %s\n", err)
		return
	}
	defer consumer.Close()

	// Listen for errors
	go func() {
		for {
			err := <-consumer.Errors()
			fmt.Printf("An error came through the Errors channel: %s\n", err)
		}
	}()

	// Consume messages forever!
	handler := simpleHandler{}
	for {
		err := consumer.Consume(context.Background(), []string{topicName}, handler)
		if err != nil {
			fmt.Printf("Failed to consume: %s\n", err)
			return
		}
	}
}
