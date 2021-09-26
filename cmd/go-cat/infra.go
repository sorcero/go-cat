// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/urfave/cli/v2"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
)

// NewInfrastructureFromCliContext converts cli.Context to Infrastructure
func NewInfrastructureFromCliContext(context *cli.Context) *infrastructure.Metadata {
	infra := &infrastructure.Metadata{
		Name:           context.String("name"),
		CommitSha:      context.String("commit-sha"),
		Version:        context.String("version"),
		Cloud:          context.String("cloud"),
		CloudProjectId: context.String("cloud-project-id"),
		Subsystem:      context.String("subsystem"),
		Type:           context.String("type"),
		MonitoringLink: context.String("monitoring-link"),
		DeploymentLink: context.String("deployment-link"),
	}

	return infra
}
