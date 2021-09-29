// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package clouds

import (
	"fmt"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
)

var googleInfrastructureMetadata = Metadata{
	Id:   "GCP",
	Name: "Google Cloud Platform",
	Types: []*TypeMetadata{
		{
			Id:   "run.googleapis.com",
			Name: "Google Cloud Run",
			GetMonitoringLink: func(m infrastructure.Metadata) string {
				return fmt.Sprintf(
					"https://console.cloud.google.com/run/detail/us-east1/%s/logs?project=%s", m.Name, m.CloudProjectId)
			},
		},
	},
}
