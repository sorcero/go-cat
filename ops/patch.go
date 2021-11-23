package ops

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"gitlab.com/sorcero/community/go-cat/parser"
	"gitlab.com/sorcero/community/go-cat/storage"
)

// Patch pulls the git repository, edits/updates the infrastructure metadata
// and pushes the repository back. All processes happen within an im-memory
// git storage system to minimize moving parts
func Patch(cfg config.GlobalConfig, infra *infrastructure.Metadata, args string) error {
	repo, fs, err := storage.Clone(cfg)
	if err != nil {
		return err
	}
	logger.Info("AAA", args)
	infra.Id = args
	logger.Info("infra.id", infra.Id)
	return PatchFromStorage(cfg, repo, fs, infra)
}

// PatchFromStorage parses the provided argument storage for infrastructure
// updates the infrastructure metadata
// and pushes the repository back. All processes happen within an im-memory
// git storage system to minimize moving parts
func PatchFromStorage(cfg config.GlobalConfig, repo *git.Repository, fs billy.Filesystem, infra *infrastructure.Metadata) error {
	infraJson, err := storage.ReadInfraDb(fs)
	if err != nil {
		return err
	}

	logger.Info("AAAA", infra)

	logger.Info("Adding infrastructure")
	infraMeta, err := parser.ReadInfrastructureFromJson(infraJson)
	if err != nil {
		return err
	}
	infraMeta, _ = infraMeta.Patch(infra)

	readmeString, infraJson, err :=	 parser.InfrastructureMetaToString(infraMeta)
	if err != nil {
		panic(err)
	}

	err = updateRepository(cfg, repo, fs, readmeString, infraJson)
	if err != nil {
		return err
	}
	return nil
}

