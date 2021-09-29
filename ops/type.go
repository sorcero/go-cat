package ops

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"gitlab.com/sorcero/community/go-cat/config"
)

type GoCatContext struct {
	Repo    *git.Repository
	Storage billy.Filesystem
	Config  config.GlobalConfig
}
