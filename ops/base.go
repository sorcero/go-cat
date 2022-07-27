package ops

import (
	"fmt"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/logging"
	"gitlab.com/sorcero/community/go-cat/storage"
	"os"
)

var logger = logging.GetLogger()

// updateRepository uses the global configuration, and creates a README.md and infra.json in the given
// in-memory file system, writes the readmeString and infraJson into the README.md and infra.json files
// and commits the changes using global git configuration. These commits are then pushed to the
// git repository as specified in the global configuration
func updateRepository(cfg config.GlobalConfig, repo *git.Repository, fs billy.Filesystem, readmeString string, infraJson []byte) error {
	readme, _ := fs.Create("README.md")
	_, err := readme.Write([]byte(readmeString))
	if err != nil {
		panic(err)
	}

	infraDb, _ := fs.Create("infra.json")
	_, err = infraDb.Write(infraJson)
	if err != nil {
		panic(err)
	}

	w, err := repo.Worktree()
	if err != nil {
		panic(err)
	}

	_, err = w.Add("README.md")
	if err != nil {
		panic(err)
	}
	_, err = w.Add("infra.json")
	if err != nil {
		panic(err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	_, err = w.Commit(fmt.Sprintf("Updated from %s", hostname), &git.CommitOptions{})

	if err != nil {
		logger.Debug(err)
		return err
	}

	logger.Info("Updating git repository with new infrastructure")
	err = storage.Push(cfg, repo)
	if err != nil {
		logger.Debug(err)
		return err
	}
	logger.Info("git repository updated")
	return nil
}
