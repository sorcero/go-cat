// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package clouds

import (
	"fmt"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"strings"
)

var googleInfrastructureMetadata = Metadata{
	Id:   "GCP",
	Name: "Google Cloud Platform",
	Types: []*TypeMetadata{
		{
			Id:   "run.googleapis.com",
			Name: "Google Cloud Run",
			GetLoggingLink: func(m infrastructure.Metadata) string {
				return fmt.Sprintf(
					"https://console.cloud.google.com/run/detail/us-east1/%s/logs?project=%s", m.Name, m.CloudProjectId)
			},
			GetMonitoringLink: func(m infrastructure.Metadata) string {
				return fmt.Sprintf(
					"https://console.cloud.google.com/run/detail/us-east1/%s/metrics?project=%s", m.Name, m.CloudProjectId)
			},
		},
		{
			Id:   "redis.googleapis.com",
			Name: "Memorystore (Redis)",
			GetLoggingLink: func(m infrastructure.Metadata) string {
				return ""
			},
			GetMonitoringLink: func(m infrastructure.Metadata) string {
				return fmt.Sprintf(
					"https://console.cloud.google.com/memorystore/redis/locations/us-east1/instances/%s/details?project=%s", m.Name, m.CloudProjectId)
			},
		},
		{
			Id:   "compute.googleapis.com",
			Name: "Google Compute Instance",
			GetLoggingLink: func(m infrastructure.Metadata) string {
				return ""
			},
			GetMonitoringLink: func(m infrastructure.Metadata) string {
				return ""
			},
		},
		{
			Id:   "container.googleapis.com/apps/v1",
			Name: "GKE Deployment",
			GetLoggingLink: func(m infrastructure.Metadata) string {
				name := m.Name
				overrideName, ok := m.Parameters["container.googleapis.com/apps/v1/metadata.name"]
				if ok && overrideName.(string) != "" {
					name = overrideName.(string)
				}

				kind := "Deployment"
				overrideKind, ok := m.Parameters["container.googleapis.com/apps/v1/kind"]
				if ok && overrideKind.(string) != "" {
					kind = overrideKind.(string)
				}
				// convert the kind to lower case because Google Cloud URLs are of lowercase
				kind = strings.ToLower(kind)

				// region, like us-east1, us-east2
				region := "us-east1"
				overrideRegion, ok := m.Parameters["container.googleapis.com/region"]
				if ok && overrideRegion.(string) != "" {
					region = overrideRegion.(string)
				}

				cluster, ok := m.Parameters["container.googleapis.com"].(string)
				if !ok {
					return ""
				}
				namespace, ok := m.Parameters["container.googleapis.com/apps/v1/namespaces"].(string)
				if !ok {
					// the default GKE namespace is default.
					namespace = "default"
				}
				return fmt.Sprintf(
					"https://console.cloud.google.com/kubernetes/%s/%s/%s/%s/%s/logs?project=%s", kind, region, cluster, namespace, name, m.CloudProjectId)
			},
			GetMonitoringLink: func(m infrastructure.Metadata) string {
				name := m.Name
				overrideName, ok := m.Parameters["container.googleapis.com/apps/v1/metadata.name"]
				if ok && overrideName != "" {
					name = overrideName.(string)
				}

				kind := "Deployment"
				overrideKind, ok := m.Parameters["container.googleapis.com/apps/v1/kind"]
				if ok && overrideKind.(string) != "" {
					kind = overrideKind.(string)
				}
				// convert the kind to lower case because Google Cloud URLs are of lowercase
				kind = strings.ToLower(kind)

				// region, like us-east1, us-east2
				region := "us-east1"
				overrideRegion, ok := m.Parameters["container.googleapis.com/region"]
				if ok && overrideRegion.(string) != "" {
					region = overrideRegion.(string)
				}

				cluster, ok := m.Parameters["container.googleapis.com"].(string)
				if !ok {
					return ""
				}
				namespace, ok := m.Parameters["container.googleapis.com/apps/v1/namespaces"].(string)
				if !ok {
					// the default GKE namespace is default.
					namespace = "default"
				}

				monitoringPage := "overview"
				if kind == "statefulset" {
					// statefulset only has a details page
					monitoringPage = "details"
				}
				return fmt.Sprintf("https://console.cloud.google.com/kubernetes/%s/%s/%s/%s/%s/%s?project=%s", kind, region, cluster, namespace, name, monitoringPage, m.CloudProjectId)
			},
		},
	},
}
