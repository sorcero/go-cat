// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/urfave/cli/v2"
	"gitlab.com/sorcero/community/go-cat/meta"
	"log"
	"os"
)

func main() {
	gitFlags := []cli.Flag{
		&cli.StringFlag{Name: "git.url", Usage: "URL to the git repository", EnvVars: []string{meta.GitUrlEnvVar}},
		&cli.StringFlag{Name: "git.username", Usage: "Username, if the Git repository requires HTTP Auth", EnvVars: []string{meta.GitUsernameEnvVar}},
		&cli.StringFlag{Name: "git.password", Usage: "Password, if the Git repository requires HTTP Auth", EnvVars: []string{meta.GitPasswordEnvVar}},
	}
	infraFlags := []cli.Flag{
		&cli.StringFlag{Name: "name", Usage: "Name of the service or endpoint"},
		&cli.StringFlag{Name: "type", Usage: "Type of infrastructure on which the service is deployed to"},
		&cli.StringFlag{Name: "commit-sha", Usage: "Deployed Commit SHA"},
		&cli.StringFlag{Name: "cloud", Usage: "Name of the cloud, example: gcp, aws"},
		&cli.StringFlag{Name: "cloud-project-id", Usage: "A unique identifier of the project / environment in which the service is deployed"},
		&cli.StringFlag{Name: "subsystem", Usage: "Name of the parent project, to which the child service belongs to"},
		&cli.StringFlag{Name: "deployment-link", Usage: "The HTTP URL or access endpoint of the API or service"},
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

				Flags: infraFlags,
			},
			{
				Name:   "push",
				Usage:  "Push changes from infrastructure queue to git",
				Action: pushInfrastructureCliContext,

				Flags: gitFlags,
			},
			{
				Name:   "cat",
				Usage:  "Read the infra.json file to stdout",
				Action: catInfrastructureCliContext,

				Flags: gitFlags,
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
		log.Fatal(err)
	}
}
