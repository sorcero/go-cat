// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package clouds

import (
	"gitlab.com/sorcero/community/go-cat/infrastructure"
)

var Meta = map[string]Metadata{
	"GCP": googleInfrastructureMetadata,
}

func getInfraMetaFromCloud(infra infrastructure.Metadata) *TypeMetadata {
	val, ok := Meta[infra.Cloud]
	if !ok {
		return nil
	}
	for i := range val.Types {
		if val.Types[i].Id == infra.Type {
			return val.Types[i]
		}
	}
	return nil
}

func GetInfraMonitoringLink(infra infrastructure.Metadata) string {
	cloudMeta := getInfraMetaFromCloud(infra)
	if cloudMeta == nil {
		return ""
	}
	return cloudMeta.GetMonitoringLink(infra)
}

func GetInfraType(infra infrastructure.Metadata) string {
	cloudMeta := getInfraMetaFromCloud(infra)
	if cloudMeta == nil {
		return ""
	}
	return cloudMeta.Name
}
