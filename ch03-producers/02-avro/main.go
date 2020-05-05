package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "avro-examples"
	app.Usage = "avro-examples"
	app.Author = "Aaron"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:   "schema-registry",
			Action: schemaRegistry,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
