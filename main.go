package main

import (
	"os"
	"sort"
	"time"

	"github.com/disaster37/harbor-cli/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

func run(args []string) error {

	// Logger setting
	log.SetOutput(os.Stdout)

	// CLI settings
	app := cli.NewApp()
	app.Usage = "Extra kubernetes tools box"
	app.Version = "develop"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "config",
			Usage: "Load configuration from `FILE`",
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "url",
			Usage:    "The Harbor base URL",
			EnvVars:  []string{"HARBOR_CLI_URL"},
			Required: true,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "username",
			Usage:    "The Harbor username to connect on it",
			EnvVars:  []string{"HARBOR_CLI_USERNAME"},
			Required: true,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:     "password",
			Usage:    "The Harbor password to connect on it",
			EnvVars:  []string{"HARBOR_CLI_PASSWORD"},
			Required: true,
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:  "timeout",
			Usage: "The timeout in second",
			Value: 60 * time.Second,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:  "self-signed-certificate",
			Usage: "Don't check remote certificat",
			Value: false,
		}),
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Display debug output",
		},
		&cli.BoolFlag{
			Name:  "no-color",
			Usage: "No print color",
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:     "check-vulnerabilities",
			Usage:    "Check vulnerabilities from Harbor",
			Category: "Artifact",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "repository",
					Usage:    "The repository name",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "project",
					Usage:    "The project name",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "artifact",
					Usage:    "The artifact name",
					Required: true,
				},
				&cli.StringFlag{
					Name:  "severity",
					Usage: "The severity level not allowed",
					Value: "",
				},
				&cli.BoolFlag{
					Name:  "force-scan",
					Usage: "Run scan on harbor before to check vulnerabilities report",
					Value: true,
				},
			},
			Action: cmd.CheckScanVulnerability,
		},
		{
			Name:     "delete-artifact",
			Usage:    "Delete artifact from Harbor",
			Category: "Artifact",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "repository",
					Usage:    "The repository name",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "project",
					Usage:    "The project name",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "artifact",
					Usage:    "The artifact name",
					Required: true,
				},
			},
			Action: cmd.CheckScanVulnerability,
		},
	}

	app.Before = func(c *cli.Context) error {

		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}

		if !c.Bool("no-color") {
			formatter := new(prefixed.TextFormatter)
			formatter.FullTimestamp = true
			formatter.ForceFormatting = true
			log.SetFormatter(formatter)
		}

		if c.String("config") != "" {
			before := altsrc.InitInputSourceWithContext(app.Flags, altsrc.NewYamlSourceFromFlagFunc("config"))
			return before(c)
		}
		return nil
	}

	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(args)
	return err
}

func main() {
	err := run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
