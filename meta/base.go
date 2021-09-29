// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package meta

import "fmt"

const (
	AppName      = "go-cat"
	EnvVarPrefix = "GO_CAT"
	QueueDbName  = ".go-cat.queue.db.json"
)

var (
	// env vars

	GitUrlEnvVar      = fmt.Sprintf("%s_GIT_URL", EnvVarPrefix)
	GitUsernameEnvVar = fmt.Sprintf("%s_GIT_USERNAME", EnvVarPrefix)
	GitPasswordEnvVar = fmt.Sprintf("%s_GIT_PASSWORD", EnvVarPrefix)
)
