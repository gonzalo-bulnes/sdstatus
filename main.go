package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/freedomofpress/sdstatus/securedrop"
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

func runScan(c *http.Client, format string, onion_urls []string) {
	i := 0

	results := make([]*securedrop.Instance, 0)

	ch := make(chan status.Information)

	// For each address we are creating a goroutine
	for _, v := range onion_urls {
		url := strings.TrimSpace(v)

		if url != "" {
			go func() {
				instance := securedrop.NewInstance(url)

				s := securedrop.NewStatusChecker(c)
				err := s.Check(instance)
				if err == nil {
					instance.Available = true
				}
				fmt.Printf("%v\n", instance)
				ch <- instance
			}()
			i = i + 1
		}

	}

	// Now wait for all the results
	for {
		result := <-ch

		if result != nil {

			if format == "csv" {
				fmt.Println(result.CSV())
			}

			results = append(results, result.(*securedrop.Instance))
			i = i - 1
		}
		if i == 0 {
			break
		}
	}

	if format == "json" {
		bits, err := json.MarshalIndent(results, "", "\t")
		if err == nil {
			fmt.Println(string(bits))
		}
	} else if format == "csv" {
	} else {
		log.Fatal(fmt.Sprintf("Invalid format: %s", format))
	}
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
		onion_urls := ctx.Args()
		if len(onion_urls) == 0 {
			log.Fatal("No args provided.")
		}

		c, err := tor.NewClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot connect to the Tor proxy: %v\n", err)
			os.Exit(1)
		}

		runScan(c, format, onion_urls)
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
