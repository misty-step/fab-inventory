# fab-inventory

A lightweight Go tool for discovering and managing repository inventories. It loads a YAML configuration file that categorizes repositories by automation tier and pipeline, providing a foundation for factory automation workflows.

## What It Does

`fab-inventory` parses a YAML inventory file containing repository metadata and provides:
- **Tier-based categorization**: Classify repos as `active`, `production`, or `dormant`
- **Pipeline association**: Track which automation pipelines apply to each repository
- **Priority labeling**: Assign priority levels (e.g., `high`, `medium`, `low`, `critical`)
- **Summary reporting**: Display counts of repos grouped by tier

This tool serves as the inventory layer for factory automation—other tools can import the `inventory` package to query repos by tier or pipeline programmatically.

## Installation

```bash
go install github.com/misty-step/fab-inventory@latest
```

Or build from source:

```bash
git clone https://github.com/misty-step/fab-inventory.git
cd fab-inventory
go install .
```

## Usage

```bash
fab-inventory <inventory.yaml>
```

### Example

```bash
$ fab-inventory inventory.yaml
✅ 4 repos loaded (2 active, 1 production, 1 dormant)
```

## Configuration

`fab-inventory` uses a YAML file to define the repository inventory.

### Inventory File Format

```yaml
org: misty-step
repos:
  factory:
    tier: active
    priority: high
    pipelines:
      - pr
      - issue-to-pr
    description: "Main factory orchestration repo"

  cerberus:
    tier: active
    priority: medium
    pipelines:
      - pr
      - backlog-groom
    description: "Issue tracking bot"

  api-service:
    tier: production
    priority: critical
    pipelines:
      - pr
      - release
    description: "Production API service"

  legacy-repo:
    tier: dormant
    priority: low
    pipelines: []
    description: "Deprecated, read-only"
```

### Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `org` | string | Yes | Organization name |
| `repos` | map | Yes | Map of repository names to repo objects |
| `repos[].tier` | string | Yes | Automation tier: `active`, `production`, or `dormant` |
| `repos[].priority` | string | No | Priority level: `high`, `medium`, `low`, `critical` |
| `repos[].pipelines` | array | No | List of pipeline names that apply to this repo |
| `repos[].description` | string | No | Human-readable description |

### Tiers

- **active**: Repositories actively developed and receiving automation (PR checks, issue workflows, etc.)
- **production**: Repositories deployed to production; typically require stricter pipeline controls
- **dormant**: Repositories no longer actively maintained; minimal or no automation

## Programmatic Usage

The `inventory` package can be imported for programmatic access:

```go
import "github.com/misty-step/fab-inventory/inventory"

// Load inventory
inv, err := inventory.LoadInventory("inventory.yaml")
if err != nil {
    log.Fatal(err)
}

// Filter by tier
activeRepos := inv.ReposForTier(inventory.TierActive)

// Filter by pipeline
prRepos := inv.ReposForPipeline("pr")
```

### Available Methods

- `LoadInventory(path string) (*Inventory, error)` — Load and parse a YAML inventory file
- `(inv *Inventory) ReposForTier(tier Tier) []Repo` — Return all repos matching the given tier
- `(inv *Inventory) ReposForPipeline(pipeline string) []Repo` — Return all repos with the specified pipeline

### Constants

```go
inventory.TierActive     // "active"
inventory.TierProduction// "production"
inventory.TierDormant    // "dormant"
```

## Integration

`fab-inventory` is designed as a foundational component in factory automation:

1. **Inventory source of truth**: The YAML file serves as the single source of truth for which repositories exist and how they're categorized
2. **Pipeline selection**: Other tools (e.g., Lobster, CI/CD systems) can query the inventory to determine which pipelines should run on which repos
3. **Reporting**: Use the CLI to get a quick overview of repository distribution across tiers

Example integration pattern:

```bash
# Get all active repos for pipeline execution
inv=$(fab-inventory inventory.yaml)
for repo in $(jq -r '.repos[] | select(.tier == "active") | .name' inventory.yaml); do
    run-pipeline $repo
done
```

## Contributing

Standard Go workflow:

```bash
# Fork and clone
git clone https://github.com/misty-step/fab-inventory.git
cd fab-inventory

# Create a feature branch
git checkout -b your-feature-name

# Run tests
go test ./...

# Run the tool
go run . your-inventory.yaml
```

### Running Tests

```bash
go test -v ./...
```

## License

MIT
