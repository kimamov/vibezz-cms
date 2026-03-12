package plugin

import (
	"context"
	"sync"

	"github.com/vibezz/cms/internal/blocks"
)

var (
	blockTypesMu   sync.RWMutex
	blockEnrichers = make(map[string]BlockEnricher)
)

// BlockTypeDefinition describes a block type available in the editor.
type BlockTypeDefinition = blocks.BlockTypeDefinition

// BlockFieldDefinition describes a field in a block type
type BlockFieldDefinition = blocks.BlockFieldDefinition

// BlockEnricher is a function that enriches block data for public rendering.
// It receives the context and the raw block data map, and returns
// the enriched data (which may include additional fields like resolved URLs,
// hydrated relations, etc.).
type BlockEnricher func(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error)

// GetAllBlockTypes returns all block type definitions including layout blocks
func GetAllBlockTypes() []BlockTypeDefinition {
	// Use the new blocks package registry
	return blocks.GetBlockTypeDefinitions()
}

// GetBlockEnricher returns an enricher for a specific block type
func GetBlockEnricher(slug string) BlockEnricher {
	blockTypesMu.RLock()
	defer blockTypesMu.RUnlock()

	// Check if we have a custom enricher
	if enricher, ok := blockEnrichers[slug]; ok {
		return enricher
	}

	// Return default enricher
	return defaultBlockEnricher
}

// defaultBlockEnricher returns data as-is
func defaultBlockEnricher(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
	return data, nil
}

// registerBlockType adds a block type to the registry
func registerBlockType(def BlockTypeDefinition) {
	blocks.RegisterBlockType(def)
}

// registerBlockEnricher adds a block enricher
func registerBlockEnricher(slug string, enricher BlockEnricher) {
	blockTypesMu.Lock()
	defer blockTypesMu.Unlock()
	blockEnrichers[slug] = enricher
}

// IsContainerBlock checks if a block type can contain other blocks
func IsContainerBlock(slug string) bool {
	return blocks.IsContainerBlock(slug)
}

// GetBlockTypeDefinition returns a specific block type definition
func GetBlockTypeDefinition(slug string) (BlockTypeDefinition, bool) {
	return blocks.GetBlockTypeDefinition(slug)
}
