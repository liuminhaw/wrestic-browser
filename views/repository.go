package views

import (
	"encoding/json"
	"fmt"

	"github.com/liuminhaw/wrestic-brw/static"
)

type RepositoryConfig struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Options     []RepositoryConfigOpts `json:"options"`
}

type RepositoryConfigOpts struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Type        string `json:"type"`
	Placeholder string `json:"placeholder"`
}

func NewRepositoryConfigs() ([]RepositoryConfig, error) {
	data, err := static.FS.ReadFile("config/repository.json")
	if err != nil {
		return nil, fmt.Errorf("new repository config: %w", err)
	}

	var configs []RepositoryConfig
	err = json.Unmarshal(data, &configs)
	if err != nil {
		return nil, fmt.Errorf("new repository config: unmarshal failed: %w", err)
	}

	return configs, nil
}
