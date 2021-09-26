// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package clouds

import "gitlab.com/sorcero/community/go-cat/infrastructure"

type Metadata struct {
	Id    string          `json:"id"`
	Name  string          `json:"name"`
	Types []*TypeMetadata `json:"types"`
}

type TypeMetadata struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	GetMonitoringLink func(metadata infrastructure.Metadata) string
}
