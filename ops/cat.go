package ops

import (
	"gitlab.com/sorcero/community/go-cat/config"
	"gitlab.com/sorcero/community/go-cat/storage"
)

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
