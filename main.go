package main

import (
	"fmt"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"

	"github.com/cookpad/itacho/generator"
	"github.com/cookpad/itacho/server"
	"github.com/cookpad/itacho/xds"
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
			Name:    "generate",
			Aliases: []string{"g"},
			Usage:   "generate xDS response",
			Action:  generate,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "source, s",
					Usage: "load service definition from `FILE`",
				},
				cli.StringFlag{
					Name:  "output, o",
					Usage: "write xDS response flagment into `DIR`",
				},
				cli.StringFlag{
					Name:  "type, t",
					Usage: "xDS response type. Currently supported types are: CDS, RDS",
				},
				cli.StringFlag{
					Name:  "version, v",
					Usage: "specify version of generated xDS response (e.g. Git sha)",
				},
				cli.BoolFlag{
					Name:  "use-legacy-sds",
					Usage: "[optional] use legacy v1 SDS instead of v2 EDS. Default `false`",
				},
				cli.StringFlag{
					Name:  "eds-cluster",
					Usage: "cluster name for EDS cluster which will be configured statically in Envoy bootstrap config",
				},
			},
		},
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

func generate(ctx *cli.Context) error {
	opts := generator.Opts{}

	source := ctx.String("source")
	if len(source) < 1 {
		return buildEmptyFlagError("source")
	}
	opts.SourcePath = source

	output := ctx.String("output")
	if len(output) < 1 {
		return buildEmptyFlagError("output")
	}
	opts.OutputDir = output

	version := ctx.String("version")
	if len(version) < 1 {
		return buildEmptyFlagError("version")
	}
	opts.Version = version

	t := ctx.String("type")
	if len(t) < 1 {
		return buildEmptyFlagError("type")
	}
	if t == "CDS" {
		opts.Type = xds.CDS

		edsCluster := ctx.String("eds-cluster")
		if len(edsCluster) < 1 {
			return buildEmptyFlagError("eds-cluster")
		}
		opts.EdsCluster = edsCluster
	} else if t == "RDS" {
		opts.Type = xds.RDS
	} else {
		return buildFlagError("type", "value must be either `CDS` or `RDS`")
	}

	legacySds := ctx.Bool("use-legacy-sds")
	if legacySds {
		opts.LegacySds = true
	}

	return generator.Generate(opts)
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
