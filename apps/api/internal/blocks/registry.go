package blocks

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/vibezz/cms/internal/domain"
)

var (
	registryMu  sync.RWMutex
	blockTypes  = make(map[string]BlockTypeDefinition)
	initialized = false
)

// BlockTypeDefinition describes a block type that can be registered
type BlockTypeDefinition struct {
	Slug        string                 `json:"slug"`
	Label       string                 `json:"label"`
	Icon        string                 `json:"icon,omitempty"`
	Category    string                 `json:"category,omitempty"`
	Description string                 `json:"description,omitempty"`
	DefaultData map[string]interface{} `json:"defaultData,omitempty"`
	IsContainer bool                   `json:"isContainer,omitempty"`
	Fields      []BlockFieldDefinition `json:"fields,omitempty"`
}

// BlockFieldDefinition describes a field in a block type
type BlockFieldDefinition struct {
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Label    string      `json:"label"`
	Required bool        `json:"required,omitempty"`
	Default  interface{} `json:"default,omitempty"`
	Options  interface{} `json:"options,omitempty"`
}

// initRegistry initializes the default block types
func initRegistry() {
	registryMu.Lock()
	defer registryMu.Unlock()

	if initialized {
		return
	}

	for _, def := range getDefaultBlockTypes() {
		blockTypes[def.Slug] = def
	}
	initialized = true
}

// RegisterBlockType adds a new block type to the registry
func RegisterBlockType(def BlockTypeDefinition) {
	initRegistry()

	registryMu.Lock()
	defer registryMu.Unlock()
	blockTypes[def.Slug] = def
}

// UnregisterBlockType removes a block type from the registry
func UnregisterBlockType(slug string) {
	initRegistry()

	registryMu.Lock()
	defer registryMu.Unlock()
	delete(blockTypes, slug)
}

// getDefaultBlockTypes returns the default block type definitions
func getDefaultBlockTypes() []BlockTypeDefinition {
	return []BlockTypeDefinition{
		// Basic blocks
		{
			Slug:        "heading",
			Label:       "Heading",
			Icon:        "heading",
			Category:    "basic",
			Description: "A heading block for titles and subtitles",
			DefaultData: map[string]interface{}{
				"text":  "",
				"level": 2,
			},
			Fields: []BlockFieldDefinition{
				{Name: "text", Type: "text", Label: "Text", Required: true},
				{Name: "level", Type: "number", Label: "Level", Default: 2},
			},
		},
		{
			Slug:        "text",
			Label:       "Text",
			Icon:        "text",
			Category:    "basic",
			Description: "A rich text paragraph block",
			DefaultData: map[string]interface{}{
				"content": "",
			},
			Fields: []BlockFieldDefinition{
				{Name: "content", Type: "richtext", Label: "Content", Required: true},
			},
		},
		{
			Slug:        "image",
			Label:       "Image",
			Icon:        "image",
			Category:    "media",
			Description: "An image block with caption",
			DefaultData: map[string]interface{}{
				"mediaId":     "",
				"caption":     "",
				"alt":         "",
				"aspectRatio": "auto",
			},
			Fields: []BlockFieldDefinition{
				{Name: "mediaId", Type: "media", Label: "Image", Required: true},
				{Name: "caption", Type: "text", Label: "Caption"},
				{Name: "alt", Type: "text", Label: "Alt Text"},
				{Name: "aspectRatio", Type: "select", Label: "Aspect Ratio", Default: "auto", Options: []string{"auto", "16:9", "4:3", "1:1", "3:2"}},
			},
		},
		{
			Slug:        "quote",
			Label:       "Quote",
			Icon:        "quote",
			Category:    "basic",
			Description: "A blockquote with attribution",
			DefaultData: map[string]interface{}{
				"text":        "",
				"attribution": "",
			},
			Fields: []BlockFieldDefinition{
				{Name: "text", Type: "richtext", Label: "Quote Text", Required: true},
				{Name: "attribution", Type: "text", Label: "Attribution"},
			},
		},
		{
			Slug:        "code",
			Label:       "Code",
			Icon:        "code",
			Category:    "basic",
			Description: "A code block with syntax highlighting",
			DefaultData: map[string]interface{}{
				"code":     "",
				"language": "plaintext",
			},
			Fields: []BlockFieldDefinition{
				{Name: "code", Type: "code", Label: "Code", Required: true},
				{Name: "language", Type: "select", Label: "Language", Default: "plaintext", Options: []string{"plaintext", "javascript", "typescript", "go", "python", "html", "css", "sql", "json"}},
			},
		},
		{
			Slug:        "divider",
			Label:       "Divider",
			Icon:        "divider",
			Category:    "basic",
			Description: "A horizontal divider line",
			DefaultData: map[string]interface{}{},
		},

		// Layout blocks
		{
			Slug:        "container",
			Label:       "Container",
			Icon:        "container",
			Category:    "layout",
			Description: "A container that groups other blocks",
			IsContainer: true,
			DefaultData: map[string]interface{}{
				"layout":     "flex",
				"gap":        "16px",
				"padding":    "0px",
				"background": "",
				"maxWidth":   "",
			},
			Fields: []BlockFieldDefinition{
				{Name: "layout", Type: "select", Label: "Layout", Default: "flex", Options: []string{"flex", "grid", "row", "column", "stack"}},
				{Name: "gap", Type: "text", Label: "Gap", Default: "16px"},
				{Name: "padding", Type: "text", Label: "Padding", Default: "0px"},
				{Name: "background", Type: "color", Label: "Background"},
				{Name: "maxWidth", Type: "text", Label: "Max Width"},
			},
		},
		{
			Slug:        "grid",
			Label:       "Grid",
			Icon:        "grid",
			Category:    "layout",
			Description: "A responsive grid layout for organizing content",
			IsContainer: true,
			DefaultData: map[string]interface{}{
				"columns":        3,
				"columnGap":      "16px",
				"rowGap":         "16px",
				"mobileColumns":  1,
				"tabletColumns":  2,
				"desktopColumns": 3,
			},
			Fields: []BlockFieldDefinition{
				{Name: "columns", Type: "number", Label: "Columns", Default: 3},
				{Name: "columnGap", Type: "text", Label: "Column Gap", Default: "16px"},
				{Name: "rowGap", Type: "text", Label: "Row Gap", Default: "16px"},
				{Name: "mobileColumns", Type: "number", Label: "Mobile Columns", Default: 1},
				{Name: "tabletColumns", Type: "number", Label: "Tablet Columns", Default: 2},
				{Name: "desktopColumns", Type: "number", Label: "Desktop Columns", Default: 3},
			},
		},
		{
			Slug:        "section",
			Label:       "Section",
			Icon:        "section",
			Category:    "layout",
			Description: "A full-width section with background",
			IsContainer: true,
			DefaultData: map[string]interface{}{
				"background":      "",
				"paddingY":        "64px",
				"paddingX":        "0px",
				"minHeight":       "",
				"backgroundImage": "",
				"backgroundSize":  "cover",
			},
			Fields: []BlockFieldDefinition{
				{Name: "background", Type: "color", Label: "Background Color"},
				{Name: "paddingY", Type: "text", Label: "Vertical Padding", Default: "64px"},
				{Name: "paddingX", Type: "text", Label: "Horizontal Padding", Default: "0px"},
				{Name: "minHeight", Type: "text", Label: "Min Height"},
				{Name: "backgroundImage", Type: "media", Label: "Background Image"},
				{Name: "backgroundSize", Type: "select", Label: "Background Size", Default: "cover", Options: []string{"cover", "contain", "auto"}},
			},
		},
		{
			Slug:        "column",
			Label:       "Column",
			Icon:        "column",
			Category:    "layout",
			Description: "A single column within a grid or row",
			IsContainer: true,
			DefaultData: map[string]interface{}{
				"width":   "",
				"padding": "0px",
			},
			Fields: []BlockFieldDefinition{
				{Name: "width", Type: "text", Label: "Width (e.g., 1/2, 1/3, 100%)"},
				{Name: "padding", Type: "text", Label: "Padding", Default: "0px"},
			},
		},
	}
}

// GetBlockTypeDefinitions returns all available block type definitions
func GetBlockTypeDefinitions() []BlockTypeDefinition {
	initRegistry()

	registryMu.RLock()
	defer registryMu.RUnlock()

	defs := make([]BlockTypeDefinition, 0, len(blockTypes))
	for _, def := range blockTypes {
		defs = append(defs, def)
	}
	return defs
}

// GetBlockTypeDefinition returns a specific block type definition
func GetBlockTypeDefinition(slug string) (BlockTypeDefinition, bool) {
	initRegistry()

	registryMu.RLock()
	defer registryMu.RUnlock()

	def, ok := blockTypes[slug]
	return def, ok
}

// IsContainerBlock returns true if the block type can contain other blocks
func IsContainerBlock(slug string) bool {
	def, ok := GetBlockTypeDefinition(slug)
	return ok && def.IsContainer
}

// BlockJSON represents a block in JSON format for serialization
type BlockJSON struct {
	ID       uuid.UUID              `json:"id"`
	Type     string                 `json:"type"`
	EntryID  uuid.UUID              `json:"entryId,omitempty"`
	ParentID *uuid.UUID             `json:"parentId,omitempty"`
	Sequence int                    `json:"sequence"`
	Data     map[string]interface{} `json:"data"`
	Style    domain.BlockStyle      `json:"style,omitempty"`
	Children []BlockJSON            `json:"children,omitempty"`
}

// SerializeBlock converts a BlockInterface to JSON
func SerializeBlock(block domain.BlockInterface) (BlockJSON, error) {
	jsonBlock := BlockJSON{
		ID:       block.ID(),
		Type:     block.BlockType(),
		EntryID:  block.EntryID(),
		ParentID: block.ParentID(),
		Sequence: block.Sequence(),
		Data:     block.BlockData(),
		Style:    block.Style(),
	}

	// Serialize children
	children := block.Children()
	if len(children) > 0 {
		jsonBlock.Children = make([]BlockJSON, len(children))
		for i, child := range children {
			childJSON, err := SerializeBlock(child)
			if err != nil {
				return BlockJSON{}, err
			}
			jsonBlock.Children[i] = childJSON
		}
	}

	return jsonBlock, nil
}

// DeserializeBlock converts JSON to BlockInterface
func DeserializeBlock(jsonBlock BlockJSON) (domain.BlockInterface, error) {
	var block domain.BlockInterface

	switch jsonBlock.Type {
	case "container":
		cb := domain.NewContainerBlock()
		if data, ok := jsonBlock.Data["layout"]; ok {
			if layout, ok := data.(string); ok {
				cb.SetLayout(layout)
			}
		}
		if data, ok := jsonBlock.Data["gap"]; ok {
			if gap, ok := data.(string); ok {
				cb.SetGap(gap)
			}
		}
		if data, ok := jsonBlock.Data["padding"]; ok {
			if padding, ok := data.(string); ok {
				cb.SetPadding(padding)
			}
		}
		if data, ok := jsonBlock.Data["background"]; ok {
			if background, ok := data.(string); ok {
				cb.SetBackground(background)
			}
		}
		if data, ok := jsonBlock.Data["maxWidth"]; ok {
			if maxWidth, ok := data.(string); ok {
				cb.SetMaxWidth(maxWidth)
			}
		}
		block = cb

	case "grid":
		gb := domain.NewGridBlock()
		if data, ok := jsonBlock.Data["columns"]; ok {
			if columns, ok := data.(float64); ok {
				gb.SetColumns(int(columns))
			}
		}
		if data, ok := jsonBlock.Data["columnGap"]; ok {
			if gap, ok := data.(string); ok {
				gb.SetColumnGap(gap)
			}
		}
		if data, ok := jsonBlock.Data["rowGap"]; ok {
			if gap, ok := data.(string); ok {
				gb.SetRowGap(gap)
			}
		}
		if data, ok := jsonBlock.Data["mobileColumns"]; ok {
			if cols, ok := data.(float64); ok {
				responsiveCols := gb.ResponsiveColumns()
				responsiveCols.Mobile = int(cols)
				gb.SetResponsiveColumns(responsiveCols)
			}
		}
		if data, ok := jsonBlock.Data["tabletColumns"]; ok {
			if cols, ok := data.(float64); ok {
				responsiveCols := gb.ResponsiveColumns()
				responsiveCols.Tablet = int(cols)
				gb.SetResponsiveColumns(responsiveCols)
			}
		}
		if data, ok := jsonBlock.Data["desktopColumns"]; ok {
			if cols, ok := data.(float64); ok {
				responsiveCols := gb.ResponsiveColumns()
				responsiveCols.Desktop = int(cols)
				gb.SetResponsiveColumns(responsiveCols)
			}
		}
		block = gb

	default:
		b := domain.NewBlock()
		b.SetBlockType(jsonBlock.Type)
		block = b
	}

	// Set common properties
	block.SetID(jsonBlock.ID)
	block.SetEntryID(jsonBlock.EntryID)
	if jsonBlock.ParentID != nil {
		block.SetParentID(jsonBlock.ParentID)
	}
	block.SetSequence(jsonBlock.Sequence)
	block.SetBlockData(jsonBlock.Data)
	block.SetStyle(jsonBlock.Style)

	// Deserialize children
	if len(jsonBlock.Children) > 0 {
		children := make([]domain.BlockInterface, len(jsonBlock.Children))
		for i, childJSON := range jsonBlock.Children {
			child, err := DeserializeBlock(childJSON)
			if err != nil {
				return nil, err
			}
			children[i] = child
		}
		block.SetChildren(children)
	}

	return block, nil
}

// FlattenBlocks flattens a nested block structure into a flat list
func FlattenBlocks(blocks []domain.BlockInterface) []domain.BlockInterface {
	var result []domain.BlockInterface

	for _, block := range blocks {
		result = append(result, block)
		if block.IsContainer() && len(block.Children()) > 0 {
			result = append(result, FlattenBlocks(block.Children())...)
		}
	}

	return result
}

// BuildBlockTree builds a tree structure from a flat list of blocks
func BuildBlockTree(blocks []domain.BlockInterface) []domain.BlockInterface {
	// Create a map for quick lookup
	blockMap := make(map[uuid.UUID]domain.BlockInterface)
	for _, block := range blocks {
		blockMap[block.ID()] = block
	}

	// Build tree
	var rootBlocks []domain.BlockInterface
	for _, block := range blocks {
		if block.ParentID() == nil {
			rootBlocks = append(rootBlocks, block)
		} else {
			// Find parent and add as child
			if parent, ok := blockMap[*block.ParentID()]; ok {
				children := append(parent.Children(), block)
				parent.SetChildren(children)
			}
		}
	}

	return rootBlocks
}

// ValidateBlock validates a block's data against its type definition
func ValidateBlock(block domain.BlockInterface) error {
	def, ok := GetBlockTypeDefinition(block.BlockType())
	if !ok {
		return fmt.Errorf("unknown block type: %s", block.BlockType())
	}

	data := block.BlockData()
	for _, field := range def.Fields {
		if field.Required {
			if _, ok := data[field.Name]; !ok || data[field.Name] == "" || data[field.Name] == nil {
				return fmt.Errorf("required field '%s' is missing in block type '%s'", field.Name, def.Slug)
			}
		}
	}

	// Validate container-specific constraints
	if def.IsContainer && len(block.Children()) == 0 {
		// Allow empty containers, they're valid
	}

	return nil
}

// SerializeBlocks serializes a list of blocks to JSON bytes
func SerializeBlocks(blocks []domain.BlockInterface) ([]byte, error) {
	jsonBlocks := make([]BlockJSON, len(blocks))
	for i, block := range blocks {
		jsonBlock, err := SerializeBlock(block)
		if err != nil {
			return nil, err
		}
		jsonBlocks[i] = jsonBlock
	}
	return json.Marshal(jsonBlocks)
}

// DeserializeBlocks deserializes JSON bytes to a list of blocks
func DeserializeBlocks(data []byte) ([]domain.BlockInterface, error) {
	var jsonBlocks []BlockJSON
	if err := json.Unmarshal(data, &jsonBlocks); err != nil {
		return nil, err
	}

	blocks := make([]domain.BlockInterface, len(jsonBlocks))
	for i, jsonBlock := range jsonBlocks {
		block, err := DeserializeBlock(jsonBlock)
		if err != nil {
			return nil, err
		}
		blocks[i] = block
	}

	return blocks, nil
}
