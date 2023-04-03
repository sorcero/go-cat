// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/withmandala/go-log"
	"os"

	"github.com/urfave/cli/v2"
	"gitlab.com/sorcero/community/go-cat/meta"
)

var logger = log.New(os.Stdout)

func main() {

	gitFlags := []cli.Flag{
		&cli.StringFlag{Name: "git.url", Usage: "URL to the git repository", EnvVars: []string{meta.GitUrlEnvVar}},
		&cli.StringFlag{Name: "git.username", Usage: "Username, if the Git repository requires HTTP Auth", EnvVars: []string{meta.GitUsernameEnvVar}},
		&cli.StringFlag{Name: "git.password", Usage: "Password, if the Git repository requires HTTP Auth", EnvVars: []string{meta.GitPasswordEnvVar}},

		&cli.StringFlag{Name: "title", Usage: "Title of the infrastructure"},
		&cli.BoolFlag{Name: "archive", Usage: "Add archive infra.json file"},
	}
	infraFlags := []cli.Flag{
		&cli.StringFlag{Name: "name", Usage: "Name of the service or endpoint"},
		&cli.StringFlag{Name: "type", Usage: "Type of infrastructure on which the service is deployed to"},
		&cli.StringFlag{Name: "labels", Usage: "Additional key:value pairs, separated by comma"},
		&cli.StringFlag{Name: "commit-sha", Usage: "Deployed Commit SHA"},
		&cli.StringFlag{Name: "cloud", Usage: "Name of the cloud, example: gcp, aws"},
		&cli.StringFlag{Name: "cloud-project-id", Usage: "A unique identifier of the project / environment in which the service is deployed"},
		&cli.StringFlag{Name: "subsystem", Usage: "Name of the parent project, to which the child service belongs to"},
		&cli.StringFlag{Name: "deployment-link", Usage: "The HTTP URL or access endpoint of the API or service"},
		&cli.StringFlag{Name: "deployment-links", Usage: "Multiple HTTP URLs or access endpoint of the API or service, separated by comma"},
		&cli.StringFlag{Name: "logging-links", Usage: "Multiple HTTP URLs of the logging dashboard of the service, separated by comma"},
		&cli.StringFlag{Name: "monitoring-links", Usage: "Multiple HTTP URLs of the monitoring dashboard of the servide, separated by comma"},
		&cli.StringFlag{Name: "parameters", Usage: "Additional parameters"},
	}
	queueFlags := []cli.Flag{
		&cli.StringFlag{Name: "queue", Required: false, Usage: "Specifies a file path to store the queue in."},
	}
	app := &cli.App{
		Name:  "go-cat",
		Usage: "CLI tool to have an overview of Infrastructure, as an API as well as Markdown",

		Commands: []*cli.Command{
			{
				Name:   "upsert",
				Usage:  "Upsert infrastructure",
				Action: upsertInfrastructureCliContext,

				Flags: append(infraFlags, gitFlags...),
			},
			{
				Name:   "add",
				Usage:  "Add infrastructure to queue",
				Action: addInfrastructureCliContext,

				Flags: append(infraFlags, queueFlags...),
			},
			{
				Name:   "push",
				Usage:  "Push changes from infrastructure queue to git",
				Action: pushInfrastructureCliContext,

				Flags: append(gitFlags, queueFlags...),
			},
			{
				Name:   "cat",
				Usage:  "Read the infra.json file to stdout",
				Action: catInfrastructureCliContext,

				Flags: append(gitFlags,
					[]cli.Flag{
						&cli.BoolFlag{Name: "deployment-link", Aliases: []string{"d"}, Usage: "Output only the deployment link"},
						&cli.IntFlag{Name: "deployment-link-index", Aliases: []string{"i"}, Usage: "Output the ith deployment link, defaults to 0", Value: 0},
					}...),
			},
			{
				Name:   "render",
				Usage:  "Convert metadata infra.json to Markdown",
				Action: renderInfrastructureCliContext,
			},
			{
				Name:   "remove",
				Usage:  "Remove a infrastructure by id, supports wildcards, w.r.t go regexp",
				Action: removeInfrastructureCliContext,

				Flags: append(gitFlags,
					[]cli.Flag{&cli.StringFlag{Name: "id", Usage: "Path to the components, supporting wildcards"}}...),
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
