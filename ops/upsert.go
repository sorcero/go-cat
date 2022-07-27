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

// Upsert pulls the git repository, adds the infrastructure metadata
// and pushes the repository back. All processes happen within an im-memory
// git storage system to minimize moving parts
func Upsert(cfg config.GlobalConfig, infra *infrastructure.Metadata) error {
	repo, fs, err := storage.Clone(cfg)
	if err != nil {
		return err
	}
	return SafeUpsertFromStorage(cfg, repo, fs, infra)
}

func SafeUpsertFromStorage(cfg config.GlobalConfig, repo *git.Repository, fs billy.Filesystem, infra *infrastructure.Metadata) error {
	operation := func() error {
		err := UpsertFromStorage(cfg, repo, fs, infra)
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

// UpsertFromStorage  parses the provided argument storage for infrastructure
// adds the infrastructure metadata
// and pushes the repository back. All processes happen within an im-memory
// git storage system to minimize moving parts
func UpsertFromStorage(cfg config.GlobalConfig, repo *git.Repository, fs billy.Filesystem, infra *infrastructure.Metadata) error {
	infraJson, err := storage.ReadInfraDb(fs)
	if err != nil {
		return err
	}

	logger.Info("Adding infrastructure")
	infraMeta, err := parser.ReadInfrastructureFromJson(infraJson)
	if err != nil {
		logger.Debug(err)
		return err
	}
	infraMeta, _ = infraMeta.Add(infra)

	readmeString, infraJson, err := parser.InfrastructureMetaToString(infraMeta)
	if err != nil {
		panic(err)
	}

	err = updateRepository(cfg, repo, fs, readmeString, infraJson)
	if err != nil {
		logger.Debug(err)
		return err
	}
	return nil
}
