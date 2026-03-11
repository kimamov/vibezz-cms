package plugin

import "sync"

var (
	blockTypesMu   sync.RWMutex
	blockTypes     []BlockTypeDefinition
	blockEnrichers = make(map[string]BlockEnricher)
)

func registerBlockType(def BlockTypeDefinition) {
	blockTypesMu.Lock()
	defer blockTypesMu.Unlock()
	blockTypes = append(blockTypes, def)
}

func registerBlockEnricher(slug string, enricher BlockEnricher) {
	blockTypesMu.Lock()
	defer blockTypesMu.Unlock()
	blockEnrichers[slug] = enricher
}

// GetAllBlockTypes returns core block types plus plugin-registered ones.
func GetAllBlockTypes() []BlockTypeDefinition {
	blockTypesMu.RLock()
	defer blockTypesMu.RUnlock()
	core := coreBlockTypes()
	out := make([]BlockTypeDefinition, 0, len(core)+len(blockTypes))
	out = append(out, core...)
	out = append(out, blockTypes...)
	return out
}

// GetBlockEnricher returns the enricher for a block type slug, or nil.
func GetBlockEnricher(slug string) BlockEnricher {
	blockTypesMu.RLock()
	defer blockTypesMu.RUnlock()
	return blockEnrichers[slug]
}

func coreBlockTypes() []BlockTypeDefinition {
	return []BlockTypeDefinition{
		{Slug: "heading", Label: "Heading", Icon: "i-heroicons-h1", DefaultData: map[string]interface{}{"text": "", "level": float64(2)}},
		{Slug: "text", Label: "Text", Icon: "i-heroicons-document-text", DefaultData: map[string]interface{}{"content": ""}},
		{Slug: "image", Label: "Image", Icon: "i-heroicons-photo", DefaultData: map[string]interface{}{"media_id": "", "caption": "", "alt": ""}},
		{Slug: "quote", Label: "Quote", Icon: "i-heroicons-chat-bubble-bottom-center-text", DefaultData: map[string]interface{}{"text": "", "attribution": ""}},
		{Slug: "code", Label: "Code", Icon: "i-heroicons-code-bracket", DefaultData: map[string]interface{}{"code": "", "language": ""}},
		{Slug: "divider", Label: "Divider", Icon: "i-heroicons-minus", DefaultData: map[string]interface{}{}},
	}
}
