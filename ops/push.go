package ops

import (
	"encoding/json"
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"gitlab.com/sorcero/community/go-cat/internal/helpers"
	"gitlab.com/sorcero/community/go-cat/meta"
	"gitlab.com/sorcero/community/go-cat/parser"
	"gitlab.com/sorcero/community/go-cat/storage"
	"os"
)

// Push pushes all the infrastructure from queue
func Push(cfg config.GlobalConfig) error {
	queueDB := meta.QueueDbName
	return PushWithDbQueue(cfg, queueDB)
}

func PushWithDbQueue(cfg config.GlobalConfig, queueDB string) error {
	repo, fs, err := storage.Clone(cfg)
	if err != nil {
		return err
	}

	infraMetaQueue := &infrastructure.MetadataGroup{}
	if helpers.CheckFileExists(queueDB) {
		data, err := os.ReadFile(queueDB)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, infraMetaQueue)
		if err != nil {
			return err
		}
	}
	infraMetaQueue.Title = cfg.Title

	err = PushFromStorage(repo, fs, infraMetaQueue, cfg)
	if err != nil {
		logger.Fatal(err)
	}
	err = os.Remove(queueDB)
	if err != nil {
		panic(err)
	}
	return nil

}

func PushFromStorage(repo *git.Repository, fs billy.Filesystem, infraMetaQueue *infrastructure.MetadataGroup, cfg config.GlobalConfig) error {
	infraMeta := &infrastructure.MetadataGroup{}

	if !cfg.Overwrite {
		infraJson, err := storage.ReadInfraDb(fs)
		if err != nil {
			return err
		}
		err = json.Unmarshal(infraJson, infraMeta)
		if err != nil {
			return err
		}
	}

	var err error
	logger.Info("Adding infrastructure")
	for i := range infraMetaQueue.Infra {
		infra := infraMetaQueue.Infra[i]

		infraMeta, err = infraMeta.Add(infra)
		if err != nil {
			return err
		}
	}
	infraMeta.Title = cfg.Title

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

	s := infraMeta.Title
	if s == "" {
		s = "Pushed"
	} else if cfg.Archive && s != "" {
		archivePath := fmt.Sprintf("archives/%s/infra.json", s)
		infraDbArchive, _ := fs.Create(archivePath)
		_, err = infraDbArchive.Write(infraJson)
		if err != nil {
			return err
		}
		_, err = w.Add(archivePath)
		if err != nil {
			return err
		}
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

	_, err = w.Commit(fmt.Sprintf("%s from %s", s, hostname), &git.CommitOptions{})
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
