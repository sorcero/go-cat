// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package config

import "github.com/urfave/cli/v2"

type GlobalConfig struct {
	// GitRepository is the url in which the infrastructure 'state' is stored
	GitRepository string

	// GitUsername and GitPassword is optional, if the repository is private, it will be required
	GitUsername string
	GitPassword string
}

func NewGlobalConfigFromCliContext(context *cli.Context) GlobalConfig {
	return GlobalConfig{
		GitRepository: context.String("git.url"),
		GitUsername:   context.String("git.username"),
		GitPassword:   context.String("git.password"),
	}
}
