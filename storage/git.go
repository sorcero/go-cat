// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package storage

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/logging"
)

var logger = logging.GetLogger()

func authFromConfig(cfg config.GlobalConfig) transport.AuthMethod {
	var auth transport.AuthMethod
	if cfg.GitPassword == "" {
		auth = nil
	} else {
		logger.Debug("Using basic auth for cloning git repository")
		auth = &http.BasicAuth{
			Username: cfg.GitUsername,
			Password: cfg.GitPassword,
		}
	}
	return auth
}

func Clone(cfg config.GlobalConfig) (*git.Repository, billy.Filesystem, error) {
	storer := memory.NewStorage()
	fs := memfs.New()
	auth := authFromConfig(cfg)
	repo, err := git.Clone(storer, fs, &git.CloneOptions{
		URL:   cfg.GitRepository,
		Auth:  auth,
		Depth: 1,
	})
	if err != nil {
		return nil, fs, err
	}
	return repo, fs, err

}

func Push(cfg config.GlobalConfig, repo *git.Repository) error {
	auth := authFromConfig(cfg)
	err := repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
	})
	if err != nil {
		logger.Debug(err)
		return err
	}
	return nil
}
