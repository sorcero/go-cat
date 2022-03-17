// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"strings"

	"github.com/urfave/cli/v2"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
)

func parseCliParameters(p string) map[string]interface{} {
	k := strings.Split(p, ",")
	m := map[string]interface{}{}
	for i := range k {
		if k[i] == "" {
			continue
		}
		keyValue := strings.SplitN(k[i], "=", 2)
		m[strings.TrimSpace(keyValue[0])] = strings.TrimSpace(keyValue[1])
	}
	return m

}

func parseDeploymentLinks(p string) []string {
	return strings.Split(p, ",")
}

// newInfrastructureFromCliContext converts cli.Context to Infrastructure
func newInfrastructureFromCliContext(context *cli.Context) *infrastructure.Metadata {
	infra := &infrastructure.Metadata{
		Name:            context.String("name"),
		CommitSha:       context.String("commit-sha"),
		Version:         context.String("version"),
		Cloud:           context.String("cloud"),
		CloudProjectId:  context.String("cloud-project-id"),
		Subsystem:       context.String("subsystem"),
		Type:            context.String("type"),
		MonitoringLink:  context.String("monitoring-link"),
		DeploymentLink:  context.String("deployment-link"),
		DeploymentLinks: parseDeploymentLinks(context.String("deployment-links")),
		MonitoringLinks: parseDeploymentLinks(context.String("monitoring-links")),
		LoggingLinks:    parseDeploymentLinks(context.String("logging-links")),
		Parameters:      parseCliParameters(context.String("parameters")),
	}

	return infra
}
