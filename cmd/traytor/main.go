package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/codegangsta/cli"
)

func showError(c *cli.Context, message string) {
	fmt.Fprintf(os.Stderr, ">>> error: %s\n\n", message)
	cli.ShowSubcommandHelp(c)
	os.Exit(1)
}

func getArguments(c *cli.Context) (string, string) {
	if c.NArg() != 2 {
		showError(c, "render arguments must be exactly 2 (input scene file and output image file)")
	}
	arguments := []string(c.Args())
	return arguments[0], arguments[1]
}

func main() {
	app := cli.NewApp()
	app.Name = "traytor test"
	app.Usage = "every single ray misses"

	app.Commands = []cli.Command{
		{
			Name:      "render",
			Aliases:   []string{"ren", "r"},
			Usage:     "render a scene locally",
			ArgsUsage: "<scene file> <output image file>",
			Action:    runRender,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "width, x",
					Usage: "width of the output image",
					Value: 800,
				},
				cli.IntFlag{
					Name:  "height, y",
					Usage: "height of the output image",
					Value: 450,
				},
				cli.IntFlag{
					Name:  "max-jobs, j",
					Value: runtime.NumCPU(),
					Usage: "number of parallel rendering threads",
				},
			},
		},
		{
			Name:    "worker",
			Aliases: []string{"wrk", "w"},
			Usage:   "start a rendering server which takes requests from the client",
			Action:  runWorker,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "listen-address, l",
					Value: ":1234",
					Usage: "local network address (interface) to bind the server to",
				},
				cli.IntFlag{
					Name:  "max-requests, r",
					Value: runtime.NumCPU() * 2,
					Usage: "max number of parallel requests to the worker",
				},
				cli.IntFlag{
					Name:  "multisample, m",
					Value: 1,
					Usage: "number of samples to send at once to the server (reduces network load)",
				},
				cli.IntFlag{
					Name:  "max-jobs, j",
					Value: runtime.NumCPU(),
					Usage: "number of parallel rendering threads",
				},
			},
		},
		{
			Name:      "client",
			Aliases:   []string{"cli", "c"},
			Usage:     "render a scene remotely on RPC workers",
			ArgsUsage: "<scene file> <output image file>",
			Action:    runClient,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "synchronous, s",
					Usage: "workers don't wait until the end to synchronise images",
				},
				cli.StringSliceFlag{
					Name:  "worker, w",
					Usage: "address of worker to connect to - can be added multiple times",
				},
				cli.IntFlag{
					Name:  "width, x",
					Usage: "width of the output image",
					Value: 800,
				},
				cli.IntFlag{
					Name:  "height, y",
					Usage: "height of the output image",
					Value: 450,
				},
			},
		},
	}

	app.Flags = []cli.Flag{}

	app.Run(os.Args)
}
