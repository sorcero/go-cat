package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/urfave/cli/v2"
	"os"

	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"gitlab.com/sorcero/community/go-cat/ops"
	"gitlab.com/sorcero/community/go-cat/parser"
)

func catInfrastructureCliContext(context *cli.Context) error {
	infra, err := ops.Cat(config.NewGlobalConfigFromCliContext(context), context.Args().First())
	if err != nil {
		return err
	}
	if context.Bool("deployment-link") {
		if len(infra) != 1 {
			return errors.New("cannot fetch deployment link for multiple infrastructure. try again without wildcards")
		}
		if infra[0].DeploymentLink != "" {
			if context.String("env") != "" {
				fmt.Printf("%s=%s\n", context.String("env"), infra[0].DeploymentLink)
			} else {
				fmt.Println(infra[0].DeploymentLink)
			}
		} else {
			if len(infra[0].DeploymentLinks) == 0 {
				return errors.New("no deployment link found for the specified resource")
			}
			if context.String("env") != "" {
				fmt.Printf("%s=%s\n", context.String("env"), infra[0].DeploymentLinks[context.Int("deployment-link-index")])
			} else {
				fmt.Println(infra[0].DeploymentLinks[context.Int("deployment-link-index")])
			}
		}

		return nil
	}

	jsonData, err := json.Marshal(infra)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}

func pushInfrastructureCliContext(context *cli.Context) error {
	o := func() error {
		var err error
		if context.String("queue") == "" {
			err = ops.Push(config.NewGlobalConfigFromCliContext(context))
		} else {
			err = ops.PushWithDbQueue(config.NewGlobalConfigFromCliContext(context), context.String("queue"))
		}

		if err != nil {
			fmt.Println(err)
			fmt.Println("retrying...")
			return err
		}
		return nil
	}
	err := backoff.Retry(o, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))
	if err != nil {
		return err
	}
	return nil
}

func renderInfrastructureCliContext(_ *cli.Context) error {
	infraMeta := &infrastructure.MetadataGroup{}
	jsonData, err := os.ReadFile("infra.json")
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

	var err error
	if context.String("queue") == "" {
		err = ops.Add(infra)
	} else {
		err = ops.AddWithDbName(infra, context.String("queue"))
	}
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

	o := func() error {
		err := ops.Upsert(config.NewGlobalConfigFromCliContext(context), i)
		if err != nil {
			fmt.Println(err)
			fmt.Println("retrying...")
			return err
		}
		return nil
	}
	err := backoff.Retry(o, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))

	if err != nil {
		return err
	}
	return nil
}

func removeInfrastructureCliContext(context *cli.Context) error {
	id := context.String("id")
	return ops.Remove(config.NewGlobalConfigFromCliContext(context), id)
}
