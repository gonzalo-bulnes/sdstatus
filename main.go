package main

import (
	"os"

	"github.com/freedomofpress/sdstatus/cli"
)

func main() {
	app, err := cli.NewApp()
	if err != nil {
		panic(err)
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
