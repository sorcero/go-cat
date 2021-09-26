// Copyright 2021 Sorcero, Inc.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/go-git/go-git/v5"
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"gitlab.com/sorcero/community/go-cat/logging"
	"gitlab.com/sorcero/community/go-cat/meta"
	"gitlab.com/sorcero/community/go-cat/parser"
	"gitlab.com/sorcero/community/go-cat/storage"

	"os"
)

var logger = logging.GetLogger()

// Upsert pulls the git repository, adds the infrastructure metadata
// and pushes the repository back. All processes happen within an im-memory
// git storage system to minimize moving parts
func Upsert(cfg config.GlobalConfig, infra *infrastructure.Metadata) error {
	repo, fs, err := storage.Clone(cfg)
	if err != nil {
		return err
	}

	infraJson, err := storage.ReadInfraDb(fs)
	if err != nil {
		return err
	}

	logger.Info("Adding infrastructure")
	infraMeta, err := infrastructure.AddInfrastructureToMarkdown(infra, infraJson)
	if err != nil {
		return err
	}

	readmeString, infraJson, err := parser.InfrastructureMetaToString(infraMeta)
	if err != nil {
		panic(err)
	}

	readme, _ := fs.Create("README.md")
	_, err = readme.Write([]byte(readmeString))
	if err != nil {
		return err
	}

	infraDb, _ := fs.Create("infra.json")
	_, err = infraDb.Write(infraJson)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add("README.md")
	if err != nil {
		return err
	}
	_, err = w.Add("infra.json")
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	_, err = w.Commit(fmt.Sprintf("Updated from %s", hostname), &git.CommitOptions{})
	if err != nil {
		return err
	}
	logger.Info("Updating git repository with new infrastructure")
	err = storage.Push(cfg, repo)
	if err != nil {
		return err
	}
	logger.Info("git repository updated")
	return nil
}

// Add adds an infrastructure to queue, which can be subsequently pushed using
// push command
func Add(infra *infrastructure.Metadata) error {
	infraMeta := &infrastructure.MetadataGroup{}
	if CheckFileExists(meta.QueueDbName) {
		data, err := ioutil.ReadFile(meta.QueueDbName)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, infraMeta)
		if err != nil {
			return err
		}
	}
	infraMeta, err := infraMeta.Add(infra)
	if err != nil {
		return err
	}
	data, err := json.Marshal(infraMeta)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(meta.QueueDbName, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

// Cat gets the infra.json and returns it as bytes
func Cat(cfg config.GlobalConfig) ([]byte, error) {
	_, fs, err := storage.Clone(cfg)
	if err != nil {
		return nil, err
	}
	infraJson, err := storage.ReadInfraDb(fs)
	if err != nil {
		return nil, err
	}
	return infraJson, nil
}

// Push pushes all the infrastructure from queue
func Push(cfg config.GlobalConfig) error {
	repo, fs, err := storage.Clone(cfg)
	if err != nil {
		return err
	}
	infraJson, err := storage.ReadInfraDb(fs)
	if err != nil {
		return err
	}

	infraMeta := &infrastructure.MetadataGroup{}
	err = json.Unmarshal(infraJson, infraMeta)
	if err != nil {
		return err
	}

	infraMetaQueue := &infrastructure.MetadataGroup{}
	if CheckFileExists(meta.QueueDbName) {
		data, err := ioutil.ReadFile(meta.QueueDbName)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, infraMetaQueue)
		if err != nil {
			return err
		}
	}

	logger.Info("Adding infrastructure")
	for i := range infraMetaQueue.Infra {
		infra := infraMetaQueue.Infra[i]

		infraMeta, err = infraMeta.Add(infra)
		if err != nil {
			return err
		}
	}

	readmeString, infraJson, err := parser.InfrastructureMetaToString(infraMeta)
	if err != nil {
		return err
	}
	readme, _ := fs.Create("README.md")
	_, err = readme.Write([]byte(readmeString))
	if err != nil {
		return err
	}

	infraDb, _ := fs.Create("infra.json")
	_, err = infraDb.Write(infraJson)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add("README.md")
	if err != nil {
		return err
	}
	_, err = w.Add("infra.json")
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	_, err = w.Commit(fmt.Sprintf("Updated from %s", hostname), &git.CommitOptions{})
	if err != nil {
		return err
	}
	logger.Info("Updating git repository with new infrastructure")
	err = storage.Push(cfg, repo)
	if err != nil {
		return err
	}
	logger.Info("git repository updated")
	err = os.Remove(meta.QueueDbName)
	if err != nil {
		return err
	}
	return nil
}

func Remove(cfg config.GlobalConfig, id string) error {
	repo, fs, err := storage.Clone(cfg)
	if err != nil {
		return err
	}

	infraJson, err := storage.ReadInfraDb(fs)
	if err != nil {
		return err
	}

	logger.Info("Adding infrastructure")

	infraMeta, err := infrastructure.RemoveInfrastructureToMarkdown(id, infraJson)
	if err != nil {
		return err
	}
	readmeString, infraJson, err := parser.InfrastructureMetaToString(infraMeta)
	if err != nil {
		panic(err)
	}

	if err != nil {
		return err
	}
	readme, _ := fs.Create("README.md")
	_, err = readme.Write([]byte(readmeString))
	if err != nil {
		return err
	}

	infraDb, _ := fs.Create("infra.json")
	_, err = infraDb.Write(infraJson)
	if err != nil {
		return err
	}

	w, err := repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add("README.md")
	if err != nil {
		return err
	}
	_, err = w.Add("infra.json")
	if err != nil {
		return err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	_, err = w.Commit(fmt.Sprintf("Updated from %s", hostname), &git.CommitOptions{})
	if err != nil {
		return err
	}
	logger.Info("Updating git repository with new infrastructure")
	err = storage.Push(cfg, repo)
	if err != nil {
		return err
	}
	logger.Info("git repository updated")
	return nil

}
