package ops

import (
	"github.com/cenkalti/backoff/v4"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"gitlab.com/sorcero/community/go-cat/parser"
	"gitlab.com/sorcero/community/go-cat/storage"
)

func Remove(cfg config.GlobalConfig, id string) error {
	repo, fs, err := storage.Clone(cfg)
	if err != nil {
		return err
	}

	return SafeRemoveFromStorage(cfg, repo, fs, id)
}

func SafeRemoveFromStorage(cfg config.GlobalConfig, repo *git.Repository, fs billy.Filesystem, id string) error {
	operation := func() error {
		err := RemoveFromStorage(cfg, repo, fs, id)
		if err != nil {
			var errClone error
			repo, fs, errClone = storage.Clone(cfg)
			if errClone != nil {
				panic(errClone)
			}
		}
		return err
	}
	return backoff.Retry(operation, backoff.NewExponentialBackOff())
}

func RemoveFromStorage(cfg config.GlobalConfig, repo *git.Repository, fs billy.Filesystem, id string) error {
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

	return updateRepository(cfg, repo, fs, readmeString, infraJson)
}
