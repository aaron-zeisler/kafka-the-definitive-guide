package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/urfave/cli"
)

// The confluent library supports a regex for the topic name, so the consumer
//  can subscribe to multiple topics using a wildcard.
func topicWildcardTest(c *cli.Context) {
	broker := "localhost:9092"
	topicName := "^test*"
	consumerGroupID := "wildcard-consumer"

	// Create a consumer and add it to a consumer group
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		// Avoid connecting to IPv6 brokers:
		// This is needed for the ErrAllBrokersDown show-case below
		// when using localhost brokers on OSX, since the OSX resolver
		// will return the IPv6 addresses first.
		// You typically don't need to specify this configuration property.
		"broker.address.family": "v4",
		"group.id":              consumerGroupID,
		"session.timeout.ms":    6000,
		"auto.offset.reset":     "earliest",
	})
	if err != nil {
		fmt.Printf("failed to create consumer: %s\n", err)
		return
	}
	defer consumer.Close()

	// Subscribe
	if err := consumer.SubscribeTopics([]string{topicName}, nil); err != nil {
		fmt.Printf("failed to subscribe: %s", err)
		return
	}

	// Create a channel for shutdown
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// Consume messages forever!
	run := true
	for run == true {
		select {
		case sig := <-sigchan:
			fmt.Printf("Caught signal %v: terminating\n", sig)
			run = false
		default:
			ev := consumer.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
				if e.Headers != nil {
					fmt.Printf("%% Headers: %v\n", e.Headers)
				}
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				// But in this example we choose to terminate
				// the application if all brokers are down.
				fmt.Printf("%% Error: %v: %v\n", e.Code(), e)
				if e.Code() == kafka.ErrAllBrokersDown {
					run = false
				}
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}
}
