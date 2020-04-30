package main

import (
	"fmt"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/cookpad/itacho/server"
)

const progName = "itacho"

// Embed at build time
var version string

func main() {
	log.SetLevel(log.DebugLevel)

	app := cli.NewApp()
	app.Name = progName
	app.Usage = "itacho to manange and operate envoy based service mesh"
	app.Version = version
	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "run xDS POST-GET convert proxy",
			Action:  serverCmd,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "addr, a",
					Usage:  "[optional] bind address like `[::]` or `0.0.0.0`",
					EnvVar: "BIND_ADDR",
				},
				cli.UintFlag{
					Name:   "port, p",
					Usage:  "[optional] bind port like `8080`.",
					EnvVar: "BIND_PORT",
				},
				cli.StringFlag{
					Name:   "object-storage-endpoint-url, e",
					Usage:  "an endpoint URL of xDS response storage",
					EnvVar: "OBJECT_STORAGE_ENDPOINT_URL",
				},
			},
		},
	}

	sort.Sort(cli.CommandsByName(app.Commands))
	sort.Sort(cli.FlagsByName(app.Flags))
	for _, c := range app.Commands {
		sort.Sort(cli.FlagsByName(c.Flags))
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func serverCmd(ctx *cli.Context) error {
	opts := server.Opts{}

	// addr can be empty
	opts.BindAddr = ctx.String("addr")
	// port can be 0
	opts.BindPort = ctx.Uint("port")

	endpoint := ctx.String("object-storage-endpoint-url")
	if len(endpoint) < 1 {
		return buildEmptyFlagError("object-storage-endpoint-url")
	}
	opts.ObjectStorageEndpoint = endpoint

	return server.Start(opts)
}

func buildEmptyFlagError(name string) error {
	return fmt.Errorf(`flag "%s" cannot be empty`, name)
}

func buildFlagError(name string, reason string) error {
	return fmt.Errorf(`Invalid flag value given for "%s": %s`, name, reason)
}
