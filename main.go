package main

import (
	"fmt"
	"os"

	"github.com/misty-step/fab-inventory/inventory"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <inventory.yaml>\n", os.Args[0])
		os.Exit(1)
	}

	path := os.Args[1]
	inv, err := inventory.LoadInventory(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to load inventory: %v\n", err)
		os.Exit(1)
	}

	// Count by tier
	active := inv.ReposForTier(inventory.TierActive)
	production := inv.ReposForTier(inventory.TierProduction)
	dormant := inv.ReposForTier(inventory.TierDormant)

	total := len(inv.Repos)
	fmt.Printf("✅ %d repos loaded (%d active, %d production, %d dormant)\n",
		total, len(active), len(production), len(dormant))
}