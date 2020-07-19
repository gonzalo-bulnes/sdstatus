package cli

import (
	"fmt"

	http "github.com/freedomofpress/sdstatus/http"
	"github.com/freedomofpress/sdstatus/tor"
	"github.com/urfave/cli"
)

// App is command-line application instance.
type App struct {
	*cli.App
	scanner *http.Scanner
	scan    func(format string, onionURLs []string) error
}

// NewApp returns an application instance ready to use.
func NewApp() (app *App, err error) {
	app = &App{
		App: cli.NewApp(),
	}
	err = app.init()

	return
}

func (app *App) init() error {
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

	client, err := tor.NewClient()
	if err != nil {
		return fmt.Errorf("error connecting to the Tor proxy: %w", err)
	}

	app.scanner = http.NewScanner(client)
	app.scan = app.scanner.Scan

	app.Action = func(ctx *cli.Context) error {
		onionURLs := ctx.Args()
		if len(onionURLs) == 0 {
			return fmt.Errorf("please provide at least one onion URL")
		}

		err = app.scan(format, onionURLs)
		if err != nil {
			return fmt.Errorf("scanning error: %w", err)
		}

		return nil
	}

	return nil
}
