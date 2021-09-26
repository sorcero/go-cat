// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package logging

import (
	"github.com/withmandala/go-log"

	"os"
)

var logger = log.New(os.Stdout)

func GetLogger() *log.Logger {
	return logger
}
