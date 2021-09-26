package ops

import (
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

	infraJson, err := storage.ReadInfraDb(fs)
	if err != nil {
		return err
	}

	logger.Info("Adding infrastructure")
	infraMeta, err := parser.ReadInfrastructureFromJson(infraJson)
	if err != nil {
		return err
	}
	infraMeta, _ = infraMeta.Add(infra)

	readmeString, infraJson, err := parser.InfrastructureMetaToString(infraMeta)
	if err != nil {
		panic(err)
	}

	err = updateRepository(cfg, repo, fs, readmeString, infraJson)
	if err != nil {
		return err
	}
	return nil
}
