// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package infrastructure

import (
	"fmt"
	"time"
)

type Metadata struct {
	// Id is automatically generated. This is the unique key used to distinguish between infrastructures
	// the Id is calculated like $cloud/$cloud_project_id/$subsystem/$name
	Id string `json:"id"`

	// Name of the component, specifically used for uniquely identifying component.
	Name string `json:"name"`

	// CommitSha, if it uses Continuous Integration to create infrastructure, you can specify which
	// commit sha triggered the build
	CommitSha string `json:"commit_sha,omitempty"`

	// Version specifies the deployed version, or branch of git repository, if it uses CI/CD
	Version string `json:"version"`

	// Branch specifies the branch of the git repository, if it uses CI/CD
	Branch string `json:"branch,omitempty"`

	// TerraformVersion specifies the version of Terraform used to create the infrastructure
	// this is only used if terraform provider is used to create the deployment
	TerraformVersion string `json:"terraform_version,omitempty"`

	// Cloud specifies the cloud provider to which it was deployed. It accepts any string which is apt
	// for your use case. Examples include GCP, AWS, Self-hosted, EU-cloud... Cloud is case-sensitive.
	Cloud string `json:"cloud"`

	// CloudProjectId is optional. It can be specified as
	CloudProjectId string `json:"cloud_project_id,omitempty"`

	// Subsystem can be considered as group of components. Can also be considered as the parent of Name
	Subsystem string `json:"subsystem,omitempty"`

	// DeployedOn is the time when which the infrastructure was added to the infra.json
	DeployedOn time.Time `json:"deployed_on"`

	// Type is the type of infrastructure on which it was deployed on, eg: run.googleapis.com, kubernetes
	Type string `json:"infra_type,omitempty"`

	// Parameters is the additional optional parameters
	Parameters map[string]interface{} `json:"parameters,omitempty"`

	// MonitoringLink helps to specify the link to monitoring, for example prometheus dashboard, etc.
	// Deprecated, use MonitoringLinks instead
	MonitoringLink string `json:"monitoring_link,omitempty"`

	// MonitoringLinks helps to specify the link to monitoring, for example prometheus dashboard, etc.
	MonitoringLinks []string `json:"monitoring_links,omitempty"`

	// LoggingLinks helps to specify links to logging dashboards, for example grafana, etc.
	LoggingLinks []string `json:"logging_links,omitempty"`

	// DeploymentLinks specifies the link to deployment, if it is HTTP API endpoint. Optional.
	// Deprecated, use DeploymentLinks instead
	DeploymentLink string `json:"deployment_link,omitempty"`

	// DeploymentLinks specifies the link to deployment, if it is HTTP API endpoint. Optional.
	DeploymentLinks []string `json:"deployment_links,omitempty"`

	// Labels are key value pairs to store additional information
	Labels map[string]string `json:"labels,omitempty"`
}

// GetId returns a unique identification id of the infrastructure
func (i *Metadata) GetId() string {
	// cloud/cloud-project-id/subsystem/component
	id := fmt.Sprintf("%s/%s/%s/%s", i.Cloud, i.CloudProjectId, i.Subsystem, i.Name)
	i.Id = id
	return i.Id
}

// MetadataGroup is the top level infrastructure data, including the
// time the entire file was updated, etc.
type MetadataGroup struct {
	Version   string      `json:"version"`
	Title     string      `json:"title"`
	UpdatedAt time.Time   `json:"updated_at"`
	Infra     []*Metadata `json:"infra"`
}
