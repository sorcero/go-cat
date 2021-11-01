// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"encoding/json"
	"fmt"
	"github.com/atsushinee/go-markdown-generator/doc"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	infraclouds "gitlab.com/sorcero/community/go-cat/infrastructure/clouds"
)

// InfrastructureMetaToString converts the metadata json (by default, infra.json) to a Markdown file
// with a table
// | Component | Subsystem | Cloud | Project ID | SHA | Deployed On | API Endpoint |
// +-----------+-----------+-------+------------+-----+-------------+--------------+
func InfrastructureMetaToString(infraMeta *infrastructure.MetadataGroup) (string, []byte, error) {
	jsonData, err := json.MarshalIndent(infraMeta, "", "  ")
	if err != nil {
		return "", nil, err
	}

	book := doc.NewMarkDown()
	book.WriteTitle("Infrastructure", 1).WriteLines(2)
	book.Write(fmt.Sprintf("Last updated on %s", infraMeta.UpdatedAt)).WriteLines(2)

	clouds := map[string]map[string][]*infrastructure.Metadata{}
	for i := range infraMeta.Infra {
		infra := infraMeta.Infra[i]
		cloud := infra.Cloud
		projectId := infra.CloudProjectId
		if _, ok := clouds[cloud]; ok {
			if _, ok := clouds[cloud][projectId]; ok {
				clouds[cloud][projectId] = append(clouds[cloud][projectId], infra)
			} else {
				clouds[cloud][projectId] = []*infrastructure.Metadata{infra}
			}
		} else {
			clouds[cloud] = map[string][]*infrastructure.Metadata{}
			clouds[cloud][projectId] = []*infrastructure.Metadata{infra}
		}
	}

	for cloudName, projects := range clouds {
		book.WriteTitle(cloudName, 3).WriteLines(2)

		for project, components := range projects {
			book.WriteTitle(project, 4).WriteLines(2)

			t := doc.NewTable(len(components), 7)
			t.SetTitle(0, "Component")
			t.SetTitle(1, "Subsystem")
			t.SetTitle(2, "Project ID")
			t.SetTitle(3, "SHA")
			t.SetTitle(4, "Deployed On")
			t.SetTitle(5, "Type")
			t.SetTitle(6, "API Endpoint")

			for i := range components {
				infra := components[i]
				t.SetContent(i, 0, infra.Name)
				t.SetContent(i, 1, infra.Subsystem)
				t.SetContent(i, 2, infra.CloudProjectId)
				t.SetContent(i, 3, fmt.Sprintf("`%s`", infra.CommitSha))
				t.SetContent(i, 4, infra.DeployedOn.Format("2006-01-02 15:04:05 -0700 MST"))

				// get monitoring links, and only show them if we support monitoring
				monitoringLinks := infraclouds.GetInfraCloudMonitoringLink(*infra)
				links := ""
				if monitoringLinks != "" {
					links = fmt.Sprintf("[(Logs 🔗)](%s)", monitoringLinks)
				}
				additionalMonitoringLinks := infraclouds.GetInfraAdditionalMonitoringLink(*infra)
				if additionalMonitoringLinks != "" {
					links += fmt.Sprintf("[(Logs 🔗)](%s)", additionalMonitoringLinks)
				}

				t.SetContent(i, 5, fmt.Sprintf("%s<br>%s", infraclouds.GetInfraType(*infra), links))
				t.SetContent(i, 6, infra.DeploymentLink)
			}
			book.WriteTable(t)
		}
	}
	return book.String(), jsonData, nil
}
