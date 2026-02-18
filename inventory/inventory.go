package inventory

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Tier represents the automation tier for a repository.
type Tier string

const (
	TierActive     Tier = "active"
	TierProduction Tier = "production"
	TierDormant    Tier = "dormant"
)

// Repo represents a single repository in the inventory.
type Repo struct {
	Name        string   `yaml:"name"` // filled from map key
	Tier        Tier     `yaml:"tier"`
	Priority    string   `yaml:"priority"`
	Pipelines   []string `yaml:"pipelines"`
	Description string   `yaml:"description"`
}

// Inventory represents the complete repository inventory.
type Inventory struct {
	Org   string          `yaml:"org"`
	Repos map[string]Repo `yaml:"repos"`
}

// LoadInventory reads and parses a YAML inventory file.
func LoadInventory(path string) (*Inventory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cannot read inventory file: %w", err)
	}

	var inv Inventory
	if err := yaml.Unmarshal(data, &inv); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}

	// Fill in Name from map keys
	for name, repo := range inv.Repos {
		repo.Name = name
		inv.Repos[name] = repo
	}

	return &inv, nil
}

// ReposForPipeline returns all repos that have the specified pipeline configured.
func (inv *Inventory) ReposForPipeline(pipeline string) []Repo {
	var result []Repo
	for _, repo := range inv.Repos {
		for _, p := range repo.Pipelines {
			if p == pipeline {
				result = append(result, repo)
				break
			}
		}
	}
	return result
}

// ReposForTier returns all repos that belong to the specified tier.
func (inv *Inventory) ReposForTier(tier Tier) []Repo {
	var result []Repo
	for _, repo := range inv.Repos {
		if repo.Tier == tier {
			result = append(result, repo)
		}
	}
	return result
}
