// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package infrastructure

import (
	"encoding/json"
	"gitlab.com/sorcero/community/go-cat/logging"
)

var logger = logging.GetLogger()

// RemoveInfrastructureToMarkdown receives the infra.json in bytes, and id string, removes the matching id
// It then converts the InfrastructureMetadata to Markdown.
func RemoveInfrastructureToMarkdown(id string, jsonData []byte) (*MetadataGroup, error) {
	// TODO: remove this function
	infraMeta := &MetadataGroup{}
	if jsonData != nil {
		err := json.Unmarshal(jsonData, infraMeta)
		if err != nil {
			return nil, err
		}
	}

	return infraMeta.Remove(id)
}
