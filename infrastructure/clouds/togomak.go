// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package clouds

import (
	"fmt"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"net/url"
	"time"
)

var togomakInfrastructureMetadata = Metadata{
	Id:   "togomak",
	Name: "Togomak",
	Types: []*TypeMetadata{
		{
			Id:   "togomak.srev.in/release",
			Name: "Release",
			GetLoggingLink: func(m infrastructure.Metadata) string {
				stage := m.Name
				overrideStage, ok := m.Parameters["togomak.srev.in/v1/stage.id"]
				if ok && overrideStage.(string) != "" {
					stage = overrideStage.(string)
				}
				instanceID := ""
				overrideInstanceId, ok := m.Parameters["togomak.srev.in/v1/instance.id"]
				if ok && overrideInstanceId.(string) != "" {
					instanceID = overrideInstanceId.(string)
				}
				gcl, ok := m.Parameters["togomak.srev.in/v1/logging"]
				if !ok || (ok && gcl.(string) != "google-cloud") {
					return ""
				}

				query := "jsonPayload.labels.stage = \"%s\"\nlabels.instanceId = \"%s\";timeRange=%s/%s--PT24H;"
				query = url.PathEscape(fmt.Sprintf(query, stage, instanceID, m.DeployedOn.Format(time.RFC3339Nano), m.DeployedOn.Format(time.RFC3339Nano)))

				return fmt.Sprintf("https://console.cloud.google.com/logs/query;query=%s?project=%s", query, m.CloudProjectId)
			},
			GetMonitoringLink: func(m infrastructure.Metadata) string {
				return ""
			},
		},
	},
}
