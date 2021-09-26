package ops

import (
	"encoding/json"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"gitlab.com/sorcero/community/go-cat/internal/helpers"
	"gitlab.com/sorcero/community/go-cat/meta"
	"io/ioutil"
)

// Add adds an infrastructure to queue, which can be subsequently pushed using
// push command
func Add(infra *infrastructure.Metadata) error {
	infraMeta := &infrastructure.MetadataGroup{}
	if helpers.CheckFileExists(meta.QueueDbName) {
		data, err := ioutil.ReadFile(meta.QueueDbName)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, infraMeta)
		if err != nil {
			return err
		}
	}
	infraMeta, err := infraMeta.Add(infra)
	if err != nil {
		return err
	}
	data, err := json.Marshal(infraMeta)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(meta.QueueDbName, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}
