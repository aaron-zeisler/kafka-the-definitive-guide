package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Simple producer examples"
	app.Usage = "Simple producer examples"
	app.Author = "Aaron"
	app.Version = "1.0.0"

	app.Commands = []cli.Command{
		{
			Name:   "fire-and-forget",
			Action: fireAndForget,
		},
		{
			Name:   "sync",
			Action: synchronous,
		},
		{
			Name:   "async",
			Action: asynchronous,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
