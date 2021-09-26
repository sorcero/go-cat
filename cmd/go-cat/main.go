// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"gitlab.com/sorcero/community/go-cat/meta"
	"gitlab.com/sorcero/community/go-cat/parser"
	"io/ioutil"
	"log"
	"os"
	"regexp"
)

func catInfrastructureCliContext(context *cli.Context) error {
	infraJson, err := Cat(config.NewGlobalConfigFromCliContext(context))
	if err != nil {
		return err
	}

	infraMeta := &infrastructure.MetadataGroup{}
	err = json.Unmarshal(infraJson, infraMeta)
	if err != nil {
		return err
	}
	var data []*infrastructure.Metadata
	args := context.Args().First()
	if args == "" {
		data = infraMeta.Infra
	} else {
		infras := infraMeta.Infra
		idRegex, err := regexp.Compile(args)
		if err != nil {
			return err
		}

		for i := range infras {
			if idRegex.MatchString(infras[i].GetId()) {
				data = append(data, infras[i])
			}
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}

func pushInfrastructureCliContext(context *cli.Context) error {
	// cloud -> cloud-project-id -> subsystem -> name
	err := Push(config.NewGlobalConfigFromCliContext(context))
	if err != nil {
		return err
	}
	return nil
}

func renderInfrastructureCliContext(_ *cli.Context) error {
	infraMeta := &infrastructure.MetadataGroup{}
	jsonData, err := ioutil.ReadFile("infra.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, infraMeta)
	if err != nil {
		return err
	}
	readme, _, err := parser.InfrastructureMetaToString(infraMeta)
	if err != nil {
		return err
	}
	fmt.Println(readme)
	return nil
}

func addInfrastructureCliContext(context *cli.Context) error {
	// cloud -> cloud-project-id -> subsystem -> name
	infra := NewInfrastructureFromCliContext(context)
	infra.GetId()

	err := Add(infra)
	if err != nil {
		return err
	}
	return nil
}

func upsertInfrastructureCliContext(context *cli.Context) error {
	//
	// cloud -> cloud-project-id -> subsystem -> name
	i := NewInfrastructureFromCliContext(context)
	i.GetId()

	err := Upsert(config.NewGlobalConfigFromCliContext(context), i)
	if err != nil {
		return err
	}
	return nil
}

func removeInfrastructureCliContext(context *cli.Context) error {
	id := context.String("id")
	return Remove(config.NewGlobalConfigFromCliContext(context), id)
}

func main() {
	gitFlags := []cli.Flag{
		&cli.StringFlag{Name: "git.url", Usage: "URL to the git repository", EnvVars: []string{fmt.Sprintf("%s_GIT_URL", meta.EnvVarPrefix)}},
		&cli.StringFlag{Name: "git.username", Usage: "Username, if the Git repository requires HTTP Auth", EnvVars: []string{fmt.Sprintf("%s_GIT_USERNAME", meta.EnvVarPrefix)}},
		&cli.StringFlag{Name: "git.password", Usage: "Password, if the Git repository requires HTTP Auth", EnvVars: []string{fmt.Sprintf("%s_GIT_PASSWORD", meta.EnvVarPrefix)}},
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
