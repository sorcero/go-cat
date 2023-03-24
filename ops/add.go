package ops

import (
	"encoding/json"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
	"gitlab.com/sorcero/community/go-cat/internal/helpers"
	"gitlab.com/sorcero/community/go-cat/meta"
	"io/ioutil"
	"os"
)

func AddWithDbName(infra *infrastructure.Metadata, queueDb string) error {
	infraMeta := &infrastructure.MetadataGroup{}
	if helpers.CheckFileExists(queueDb) {
		data, err := ioutil.ReadFile(queueDb)
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
	err = os.WriteFile(queueDb, data, 0o644)
	if err != nil {
		return err
	}

	return nil
}

// Add adds an infrastructure to queue, which can be subsequently pushed using
// push command
func Add(infra *infrastructure.Metadata) error {
	queueDb := meta.QueueDbName
	return AddWithDbName(infra, queueDb)
}
