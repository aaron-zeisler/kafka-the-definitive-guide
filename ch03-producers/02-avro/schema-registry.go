package main

import (
	"bytes"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/elodina/go-avro"
	schemaregistry "github.com/lensesio/schema-registry"
	"github.com/urfave/cli"
)

type Customer struct {
	ID           int32  `avro:"id"`
	Name         string `avro:"name"`
	EmailAddress string `avro:"email"`
}

func schemaRegistry(c *cli.Context) {
	customer := Customer{
		ID:           12345,
		Name:         "Testy McTesterson",
		EmailAddress: "testy@example.com",
	}

	registryURL := "http://localhost:8081"
	topicName := "avro-customer-test-2"

	// Look up the schema in the registry and get the schema's definition
	schema, err := getSchemaFromRegistry(registryURL, topicName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Encode the customer struct using the schema
	encodedCustomer, err := encodeStruct(&customer, schema)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		//"key.serializer":      "io.confluent.kafka.serializers.KafkaAvroSerializer",
		//"value.serializer":    "io.confluent.kafka.serializers.KafkaAvroSerializer",
		//"schema.registry.url": "localhost:8081",
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer producer.Close()

	// Define the message
	message := kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &topicName,
		},
		Value: encodedCustomer,
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

func encodeStruct(data interface{}, rawSchema string) ([]byte, error) {
	schema, err := avro.ParseSchema(rawSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the schema: %w", err)
	}

	writer := avro.NewSpecificDatumWriter()
	// SetSchema must be called before calling Write
	writer.SetSchema(schema)

	// Create a new Buffer and Encoder to write to this Buffer
	buffer := new(bytes.Buffer)
	encoder := avro.NewBinaryEncoder(buffer)

	// Write the record
	// The 'data' struct here MUST be a pointer
	if err := writer.Write(data, encoder); err != nil {
		return nil, fmt.Errorf("failed to encode the customer using the schema: %w", err)
	}

	return buffer.Bytes(), nil
}
