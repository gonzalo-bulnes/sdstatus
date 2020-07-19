package main

import (
	"fmt"
	"log"
	"os"

	_http "github.com/freedomofpress/sdstatus/http"
	"github.com/freedomofpress/sdstatus/status"
	"github.com/freedomofpress/sdstatus/tor"
	"github.com/urfave/cli"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkStatusWithStatusChecker(checker status.Checker) error {
	return nil
}

func createApp() *cli.App {
	app := cli.NewApp()
	var format string
	app.EnableBashCompletion = true
	app.Name = "sdstatus"
	app.Version = "0.1.0"
	app.Usage = "To scan SecureDrop instances"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "format",
			Usage:       "Output scan results in `FORMAT`",
			Value:       "json",
			Destination: &format,
		},
	}
	app.Action = func(ctx *cli.Context) error {
		onionURLs := ctx.Args()
		if len(onionURLs) == 0 {
			log.Fatal("Please provide at least one onion URL.")
		}

		client, err := tor.NewClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot connect to the Tor proxy: %v\n", err)
			os.Exit(1)
		}

		s := _http.NewScanner(client)
		s.Scan(format, onionURLs)

		return nil
	}

	return app
}

func main() {
	app := createApp()
	if err := app.Run(os.Args); err != nil {
		check(err)
	}
}
