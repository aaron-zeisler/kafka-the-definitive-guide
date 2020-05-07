package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "consumer-examples"
	app.Usage = "consumer-examples"
	app.Author = "Aaron"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:   "simple-consumer",
			Action: simpleConsumer,
		},
		{
			Name:   "avro",
			Action: avroConsumer,
		},
		{
			Name:   "consume-specific-partition",
			Action: consumeSpecificPartition,
		},
		{
			Name:   "topic-wildcard",
			Action: topicWildcardTest,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
