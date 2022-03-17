// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package parser

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/atsushinee/go-markdown-generator/doc"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	infraclouds "gitlab.com/sorcero/community/go-cat/infrastructure/clouds"
	"gitlab.com/sorcero/community/go-cat/logging"
)

var logger = logging.GetLogger()

// InfrastructureMetaToString converts the metadata json (by default, infra.json) to a Markdown file
// with a table
// | Component | Subsystem | Cloud | SHA | Deployed On | API Endpoint |
// +-----------+-----------+-------+-----+-------------+--------------+
func InfrastructureMetaToString(infraMeta *infrastructure.MetadataGroup) (string, []byte, error) {
	jsonData, err := json.MarshalIndent(infraMeta, "", "  ")
	if err != nil {
		return "", nil, err
	}

	book := doc.NewMarkDown()
	book.WriteTitle("Infrastructure", 1).WriteLines(2)
	book.Write(fmt.Sprintf("Last updated on %s", infraMeta.UpdatedAt)).WriteLines(2)

	legend := doc.NewTable(3, 3)
	legend.SetTitle(0, "Icon")
	legend.SetTitle(1, "Legend")
	legend.SetTitle(1, "Description")
	// monitoring
	legend.SetContent(0, 0, "📈")
	legend.SetContent(0, 1, "Monitoring")
	legend.SetContent(0, 2, "Links to cloud monitoring services, like CPU, RAM, Requests, etc.")
	// logging
	legend.SetContent(1, 0, "📜")
	legend.SetContent(1, 1, "Logging")
	legend.SetContent(1, 2, "Links to logging dashboards for Application generated logs")
	// custom
	legend.SetContent(2, 0, "✏️")
	legend.SetContent(2, 1, "Custom")
	legend.SetContent(2, 2, "Custom links to monitoring dashboards like Grafana, which is not cloud native.")
	book.WriteTable(legend).WriteLines(2)

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

			book.Write("<details>").WriteLines(2)
			book.Write("<summary>Click to expand</summary>").WriteLines(2)

			t := doc.NewTable(len(components), 6)
			t.SetTitle(0, "Component")
			t.SetTitle(1, "Subsystem")
			t.SetTitle(2, "SHA")
			t.SetTitle(3, "Deployed On")
			t.SetTitle(4, "Type")
			t.SetTitle(5, "API Endpoints")

			for i := range components {
				infra := components[i]
				t.SetContent(i, 0, fmt.Sprintf("**%s**", infra.Name))
				t.SetContent(i, 1, infra.Subsystem)
				t.SetContent(i, 2, fmt.Sprintf("`%s`", infra.CommitSha))
				t.SetContent(i, 3, infra.DeployedOn.Format("2006-01-02 15:04:05 -0700 MST"))

				// get logging and monitoring links, and only show them if we support monitoring
				links := []string{}
				linksCount := 0
				monitoringLinks := infraclouds.GetInfraCloudMonitoringLink(*infra)
				if monitoringLinks != "" {
					linksCount++
					links = append(links, fmt.Sprintf("[\\[%d\\] 📈🔗](%s)", linksCount, monitoringLinks))
				}
				additionalMonitoringLinks := infraclouds.GetInfraAdditionalMonitoringLink(*infra)
				if additionalMonitoringLinks != "" {
					linksCount++
					links = append(links, fmt.Sprintf("[\\[%d\\] 📜🔗](%s)", linksCount, additionalMonitoringLinks))
				}
				loggingLinks := infraclouds.GetInfraCloudLoggingLink(*infra)
				if loggingLinks != "" {
					linksCount++
					links = append(links, fmt.Sprintf("[\\[%d\\] 📜🔗](%s)", linksCount, loggingLinks))
				}
				userDefinedMonitoringLinks := (*infra).MonitoringLinks
				if len(userDefinedMonitoringLinks) > 0 {
					for i := range userDefinedMonitoringLinks {
						linksCount++
						links = append(links, fmt.Sprintf("[\\[%d\\] ✏️📈🔗](%s)", linksCount, userDefinedMonitoringLinks[i]))
					}
				}

				userDefinedLoggingLinks := (*infra).LoggingLinks
				if len(userDefinedLoggingLinks) > 0 {
					for i := range userDefinedLoggingLinks {
						linksCount++
						links = append(links, fmt.Sprintf("[\\[%d\\] ✏️📜🔗](%s)", linksCount, userDefinedLoggingLinks[i]))
					}
				}

				t.SetContent(i, 4, fmt.Sprintf("%s<br>%s", infraclouds.GetInfraType(*infra), strings.Join(links, "<br>")))
				if infra.DeploymentLink != "" {
					logger.Warn("infra.DeploymentLink is deprecated and will be removed in a future version, use infra.DeploymentLinks instead.")
					t.SetContent(i, 5, infra.DeploymentLink)
				} else {
					var deploymentLinksEnumerated []string
					for link := range infra.DeploymentLinks {
						deploymentLinksEnumerated = append(deploymentLinksEnumerated, fmt.Sprintf("%d. %s", link+1, infra.DeploymentLinks[link]))
					}
					t.SetContent(i, 5, strings.Join(deploymentLinksEnumerated, "<br>"))
				}
			}
			book.WriteTable(t).WriteLines(2)
			book.Write("</details>").WriteLines(2)
		}
	}
	return book.String(), jsonData, nil
}
