// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package storage

import (
	"github.com/go-git/go-billy/v5"
	"io/ioutil"
	"os"
)

func ReadInfraDb(fs billy.Filesystem) ([]byte, error) {
	var infraJson []byte

	file, err := fs.Open("infra.json")
	if err != nil && os.IsNotExist(err) {
		file, err = fs.Create("infra.json")
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		infraJson, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
	}
	return infraJson, nil
}
