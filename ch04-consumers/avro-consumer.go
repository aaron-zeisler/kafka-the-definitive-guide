package main

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/elodina/go-avro"
	schemaregistry "github.com/lensesio/schema-registry"
	"github.com/urfave/cli"
)

// This avroHandler implements the sarama.ConsumerGroupHandler interface.
// A ConsumerGroupHandler must be passed to sarama's Consume() method.
type avroHandler struct {
	schema string
}

func (h avroHandler) Setup(session sarama.ConsumerGroupSession) error   { return nil }
func (h avroHandler) Cleanup(session sarama.ConsumerGroupSession) error { return nil }
func (h avroHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// Instantiate an empty CUstomer object
		customer := Customer{}

		// Decode the incoming message using the AVRO schema
		if err := decodeStruct(msg.Value, h.schema, &customer); err != nil {
			return fmt.Errorf("failed to decode the message into a struct: %w", err)
		}
		fmt.Printf("Received a message! %#v\n", customer)
	}
	return nil
}

type Customer struct {
	ID           int32  `avro:"id"`
	Name         string `avro:"name"`
	EmailAddress string `avro:"email"`
}

func avroConsumer(c *cli.Context) {
	broker := "localhost:9092"
	registryURL := "http://localhost:8081"
	topicName := "avro-customer-test-2"
	consumerGroupID := "avro-consumer"

	// Retrieve the schema for this topic and attach it to a handler
	schema, err := getSchemaFromRegistry(registryURL, topicName)
	if err != nil {
		fmt.Printf("failed to retrieve the AVRO schema: %s\n", err)
		return
	}

	handler := avroHandler{
		schema: schema,
	}

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
	for {
		err := consumer.Consume(context.Background(), []string{topicName}, handler)
		if err != nil {
			fmt.Printf("Failed to consume: %s\n", err)
			return
		}
	}
}

func getSchemaFromRegistry(registryURL, topicName string) (string, error) {
	schemaSubject := fmt.Sprintf("%s-value", topicName)

	registryClient, err := schemaregistry.NewClient(registryURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to the registry: %w", err)
	}

	schema, err := registryClient.GetLatestSchema(schemaSubject)
	if err != nil {
		return "", fmt.Errorf("failed to get the schema from the registry: %w", err)
	}

	return schema.Schema, nil
}

func decodeStruct(data []byte, rawSchema string, target interface{}) error {
	schema, err := avro.ParseSchema(rawSchema)
	if err != nil {
		return fmt.Errorf("failed to parse the schema: %w", err)
	}

	reader := avro.NewSpecificDatumReader()
	// SetSchema must be called before calling Read
	reader.SetSchema(schema)

	// Create a new Decoder with a given buffer
	decoder := avro.NewBinaryDecoder(data)

	// Read data into a given record with a given Decoder.
	err = reader.Read(target, decoder)
	if err != nil {
		return fmt.Errorf("failed to decode the message into the struct: %w", err)
	}

	return nil
}
