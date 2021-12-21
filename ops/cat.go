package ops

import (
	"encoding/json"
	"regexp"

	"github.com/go-git/go-billy/v5"
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"gitlab.com/sorcero/community/go-cat/storage"
)

// Cat gets the infra.json and returns it as bytes
func Cat(cfg config.GlobalConfig, args string) ([]*infrastructure.Metadata, error) {
	_, fs, err := storage.Clone(cfg)
	if err != nil {
		return nil, err
	}

	data, err := CatFromStorage(fs, args)
	if err != nil {
		return nil, err
	}

	return data, nil

}

// CatFromStorage parses the provided argument storage for infrastructure, and returns results
func CatFromStorage(fs billy.Filesystem, args string) ([]*infrastructure.Metadata, error) {
	infraJson, err := storage.ReadInfraDb(fs)
	if err != nil {
		return nil, err
	}

	infraMeta := &infrastructure.MetadataGroup{}
	err = json.Unmarshal(infraJson, infraMeta)
	if err != nil {
		return nil, err
	}
	var data []*infrastructure.Metadata

	if args == "" {
		data = infraMeta.Infra
	} else {
		infras := infraMeta.Infra
		idRegex, err := regexp.Compile(args + "$") // match the end of the string
		if err != nil {
			return nil, err
		}

		for i := range infras {
			if idRegex.MatchString(infras[i].GetId()) {
				data = append(data, infras[i])
			}
		}
	}
	return data, nil
}
