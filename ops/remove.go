package ops

import (
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
