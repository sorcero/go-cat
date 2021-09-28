package main

import (
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"gitlab.com/sorcero/community/go-cat/ops"
	"gitlab.com/sorcero/community/go-cat/parser"
	"io/ioutil"
)

func catInfrastructureCliContext(context *cli.Context) error {
	infra, err := ops.Cat(config.NewGlobalConfigFromCliContext(context), context.Args().First())
	if err != nil {
		return err
	}

	jsonData, err := json.Marshal(infra)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

func pushInfrastructureCliContext(context *cli.Context) error {
	err := ops.Push(config.NewGlobalConfigFromCliContext(context))
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
	infra := newInfrastructureFromCliContext(context)
	infra.GetId()

	err := ops.Add(infra)
	if err != nil {
		return err
	}
	return nil
}

func upsertInfrastructureCliContext(context *cli.Context) error {
	//
	// cloud -> cloud-project-id -> subsystem -> name
	i := newInfrastructureFromCliContext(context)
	i.GetId()

	err := ops.Upsert(config.NewGlobalConfigFromCliContext(context), i)
	if err != nil {
		return err
	}
	return nil
}

func removeInfrastructureCliContext(context *cli.Context) error {
	id := context.String("id")
	return ops.Remove(config.NewGlobalConfigFromCliContext(context), id)
}
