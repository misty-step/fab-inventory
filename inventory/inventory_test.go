package inventory

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadInventory(t *testing.T) {
	// Test valid YAML loads correctly
	t.Run("valid YAML loads correctly", func(t *testing.T) {
		yamlContent := `org: misty-step
repos:
  factory:
    tier: active
    priority: high
    pipelines: [pr, issue-to-pr]
    description: "Test repo"
`
		tmpFile := filepath.Join(t.TempDir(), "inventory.yaml")
		if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
			t.Fatal(err)
		}

		inv, err := LoadInventory(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if inv.Org != "misty-step" {
			t.Errorf("expected org misty-step, got %q", inv.Org)
		}

		repo, ok := inv.Repos["factory"]
		if !ok {
			t.Fatal("expected factory repo in inventory")
		}

		if repo.Name != "factory" {
			t.Errorf("expected name factory, got %q", repo.Name)
		}
		if repo.Tier != TierActive {
			t.Errorf("expected tier active, got %v", repo.Tier)
		}
		if repo.Priority != "high" {
			t.Errorf("expected priority high, got %q", repo.Priority)
		}
	})

	// Test invalid YAML returns error
	t.Run("invalid YAML returns error", func(t *testing.T) {
		yamlContent := `org: misty-step
repos:
  factory:
    tier: active
    pipelines: [pr]
`
		tmpFile := filepath.Join(t.TempDir(), "inventory.yaml")
		if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
			t.Fatal(err)
		}

		// Test file not found
		_, err := LoadInventory("/nonexistent/path.yaml")
		if err == nil {
			t.Error("expected error for nonexistent file")
		}
	})

	// Test empty inventory is valid
	t.Run("empty inventory is valid", func(t *testing.T) {
		yamlContent := `org: misty-step
repos:
`
		tmpFile := filepath.Join(t.TempDir(), "inventory.yaml")
		if err := os.WriteFile(tmpFile, []byte(yamlContent), 0644); err != nil {
			t.Fatal(err)
		}

		inv, err := LoadInventory(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(inv.Repos) != 0 {
			t.Errorf("expected 0 repos, got %d", len(inv.Repos))
		}
	})
}

func TestReposForPipeline(t *testing.T) {
	inv := &Inventory{
		Org: "misty-step",
		Repos: map[string]Repo{
			"factory": {
				Name:      "factory",
				Tier:      TierActive,
				Pipelines: []string{"pr", "issue-to-pr"},
			},
			"cerberus": {
				Name:      "cerberus",
				Tier:      TierActive,
				Pipelines: []string{"pr", "backlog-groom"},
			},
			"dormant-repo": {
				Name:      "dormant-repo",
				Tier:      TierDormant,
				Pipelines: []string{},
			},
		},
	}

	// Test ReposForPipeline filters correctly
	t.Run("filters by pipeline correctly", func(t *testing.T) {
		repos := inv.ReposForPipeline("pr")
		if len(repos) != 2 {
			t.Errorf("expected 2 repos for pr pipeline, got %d", len(repos))
		}

		repos = inv.ReposForPipeline("issue-to-pr")
		if len(repos) != 1 {
			t.Errorf("expected 1 repo for issue-to-pr pipeline, got %d", len(repos))
		}
		if repos[0].Name != "factory" {
			t.Errorf("expected factory repo, got %q", repos[0].Name)
		}
	})

	// Test unknown pipeline returns empty results
	t.Run("unknown pipeline returns empty results", func(t *testing.T) {
		repos := inv.ReposForPipeline("nonexistent")
		if len(repos) != 0 {
			t.Errorf("expected 0 repos for unknown pipeline, got %d", len(repos))
		}
	})
}

func TestReposForTier(t *testing.T) {
	inv := &Inventory{
		Org: "misty-step",
		Repos: map[string]Repo{
			"factory": {
				Name: "factory",
				Tier: TierActive,
			},
			"cerberus": {
				Name: "cerberus",
				Tier: TierActive,
			},
			"production-repo": {
				Name: "production-repo",
				Tier: TierProduction,
			},
			"dormant-repo": {
				Name: "dormant-repo",
				Tier: TierDormant,
			},
		},
	}

	// Test ReposForTier filters correctly
	t.Run("filters by tier correctly", func(t *testing.T) {
		repos := inv.ReposForTier(TierActive)
		if len(repos) != 2 {
			t.Errorf("expected 2 repos for active tier, got %d", len(repos))
		}

		repos = inv.ReposForTier(TierProduction)
		if len(repos) != 1 {
			t.Errorf("expected 1 repo for production tier, got %d", len(repos))
		}
		if repos[0].Name != "production-repo" {
			t.Errorf("expected production-repo, got %q", repos[0].Name)
		}
	})

	// Test unknown tier returns empty results
	t.Run("unknown tier returns empty results", func(t *testing.T) {
		repos := inv.ReposForTier(Tier("unknown"))
		if len(repos) != 0 {
			t.Errorf("expected 0 repos for unknown tier, got %d", len(repos))
		}
	})
}
