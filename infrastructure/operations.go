// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package infrastructure

import (
	"regexp"
	"time"
)

// Remove removes Infrastructure from InfrastructureMetadata which matches the id
// using regex.
//      $ go-cat remove --id "aws/*"
func (infraMeta *MetadataGroup) Remove(id string) (*MetadataGroup, error) {
	infras := infraMeta.Infra
	idRegex, err := regexp.Compile(id)
	if err != nil {
		return nil, err
	}

	var newInfras []*Metadata
	for i := range infras {
		if !idRegex.MatchString(infras[i].GetId()) {
			newInfras = append(newInfras, infras[i])
		} else {
			logger.Info("Removing", infras[i].GetId())
		}
	}
	return &MetadataGroup{Infra: newInfras, Version: "1", UpdatedAt: time.Now()}, nil
}

// Add updates InfrastructureMetadata with Infrastructure, duplicate infrastructure
// is merged with each other, and the final InfrastructureMetadata is returned
func (infraMeta *MetadataGroup) Add(infra *Metadata) (*MetadataGroup, error) {
	infra.DeployedOn = time.Now()
	infras := infraMeta.Infra
	for i := range infras {
		if infras[i].GetId() == infra.GetId() {
			// we found a match of the same infra, but probably older
			infras[i] = infra
			return &MetadataGroup{Infra: infras, Version: "1", UpdatedAt: time.Now()}, nil
		}
	}
	infras = append(infras, infra)
	return &MetadataGroup{Infra: infras, Version: "1", UpdatedAt: time.Now()}, nil
}


// Patch updates InfrastructureMetadata with Infrastructure, duplicate infrastructure
// is merged with each other, and the final InfrastructureMetadata is returned
func (infraMeta *MetadataGroup) Patch(infra *Metadata) (*MetadataGroup, error) {
	infra.DeployedOn = time.Now()
	infras := infraMeta.Infra
	for i := range infras {
		newInfra := infras[i]
		logger.Info(infra.GetId(), newInfra.GetId())
		if newInfra.GetId() == infra.GetId() {
			// we found a match of the same infra, but probably older
			if infra.Parameters != nil {
				newInfra.Parameters = infra.Parameters
			}
			if infra.MonitoringLink != "" {
				newInfra.MonitoringLink = infra.MonitoringLink
			}
			if infra.DeploymentLinks != nil {
				newInfra.DeploymentLinks = infra.DeploymentLinks
			}
			if infra.DeploymentLink != "" {
				newInfra.DeploymentLink = infra.DeploymentLink
			}
			if infra.CommitSha != "" {
				newInfra.CommitSha = infra.CommitSha
			}
			if infra.Type != "" {
				newInfra.Type = infra.Type
			}
			infras[i] = newInfra
			logger.Info(newInfra)

			return &MetadataGroup{Infra: infras, Version: "1", UpdatedAt: time.Now()}, nil
		}
	}
	infras = append(infras, infra)
	return &MetadataGroup{Infra: infras, Version: "1", UpdatedAt: time.Now()}, nil
}

