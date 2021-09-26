package parser

import (
	"encoding/json"
	"gitlab.com/sorcero/community/go-cat/infrastructure"
)

// ReadInfrastructureFromJson reads json file, and marshall them into
// the infrastructure.MetadataGroup which can then be used internally within
// go-cat. If data provided is empty, an empty infrastructure metadata will be
// returned instead
func ReadInfrastructureFromJson(data []byte) (*infrastructure.MetadataGroup, error) {
	infraMeta := &infrastructure.MetadataGroup{}
	if data != nil {
		err := json.Unmarshal(data, infraMeta)
		if err != nil {
			return nil, err
		}
	}
	return infraMeta, nil
}
