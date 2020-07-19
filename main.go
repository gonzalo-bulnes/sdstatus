package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/freedomofpress/sdstatus/securedrop"
	"github.com/urfave/cli"
	"golang.org/x/net/proxy"
)

const (
	// proxyAddr points to local SOCKS proxy from Tor
	proxyAddr = "127.0.0.1:9050"
)

// Information represents data that can be serialized as CSV
type Information interface {
	CSV() string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkStatus(ch chan Information, client *http.Client, url string) {
	instance := securedrop.NewInstance(url)

	// Create the request
	req, err := instance.NewMetadataRequest()
	if err != nil {
		instance.Available = false
		ch <- instance
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		instance.Available = false
		ch <- instance
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		instance.Available = false
		ch <- instance
		return
	}

	var info securedrop.Metadata
	json.Unmarshal(body, &info)

	instance.Info = info
	instance.Available = true
	ch <- instance
}

func runScan(format string, onion_urls []string) {
	i := 0

	results := make([]securedrop.Instance, 0)
	// create a SOCKS5 dialer
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		fmt.Fprintln(os.Stderr, "can't connect to the proxy:", err)
		os.Exit(1)
	}
	// setup the http client
	httpTransport := &http.Transport{}
	c := &http.Client{Transport: httpTransport}
	// Add the dialer
	httpTransport.Dial = dialer.Dial

	ch := make(chan Information)

	// For each address we are creating a goroutine
	for _, v := range onion_urls {
		url := strings.TrimSpace(v)

		if url != "" {
			go checkStatus(ch, c, v)
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

			results = append(results, result.(securedrop.Instance))
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
	app.Action = func(c *cli.Context) error {
		onion_urls := c.Args()
		if len(onion_urls) == 0 {
			log.Fatal("No args provided.")
		}
		runScan(format, onion_urls)
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
