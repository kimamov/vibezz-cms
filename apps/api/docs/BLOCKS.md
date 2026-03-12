# Block System Documentation

## Overview

The block system provides a flexible way to build content using composable blocks. It supports both basic content blocks (heading, text, image) and layout blocks (container, grid, section) that can contain nested blocks.

## Block Types

### Basic Blocks

| Block | Slug | Description |
|-------|------|-------------|
| Heading | `heading` | H1-H6 headings with configurable level |
| Text | `text` | Rich text paragraphs |
| Image | `image` | Image with caption, alt text, and aspect ratio |
| Quote | `quote` | Blockquote with attribution |
| Code | `code` | Syntax-highlighted code blocks |
| Divider | `divider` | Horizontal divider line |

### Layout Blocks

| Block | Slug | Description |
|-------|------|-------------|
| Container | `container` | Flexible container with layout options |
| Grid | `grid` | Responsive grid layout |
| Section | `section` | Full-width section with background |
| Column | `column` | Single column within a grid/row |

## Using Blocks

### Creating a Basic Block

```go
import "github.com/vibezz/cms/internal/domain"

// Create a heading block
heading := domain.NewBlock()
heading.SetBlockType("heading")
heading.SetBlockData(map[string]interface{}{
    "text": "Hello World",
    "level": 1,
})
```

### Creating a Container Block

```go
// Create a container
container := domain.NewContainerBlock()
container.SetLayout("flex")
container.SetGap("16px")
container.SetPadding("24px")
container.SetBackground("#f5f5f5")

// Add child blocks
child1 := domain.NewBlock()
child1.SetBlockType("text")
child1.SetBlockData(map[string]interface{}{
    "content": "First paragraph",
})

child2 := domain.NewBlock()
child2.SetBlockType("text")
child2.SetBlockData(map[string]interface{}{
    "content": "Second paragraph",
})

container.SetChildren([]domain.BlockInterface{child1, child2})
```

### Creating a Grid Block

```go
// Create a responsive grid
grid := domain.NewGridBlock()
grid.SetColumns(3)
grid.SetColumnGap("16px")
grid.SetRowGap("16px")

// Configure responsive columns
grid.SetResponsiveColumns(domain.ResponsiveColumns{
    Mobile:  1,  // 1 column on mobile
    Tablet:  2,  // 2 columns on tablet
    Desktop: 3,  // 3 columns on desktop
})

// Add items to the grid
items := make([]domain.BlockInterface, 6)
for i := 0; i < 6; i++ {
    item := domain.NewBlock()
    item.SetBlockType("text")
    item.SetBlockData(map[string]interface{}{
        "content": fmt.Sprintf("Grid item %d", i+1),
    })
    items[i] = item
}
grid.SetChildren(items)
```

## Serializing Blocks

```go
import "github.com/vibezz/cms/internal/blocks"

// Serialize a single block
block := createSomeBlock()
jsonData, err := blocks.SerializeBlock(block)
if err != nil {
    log.Fatal(err)
}

// Serialize multiple blocks
blockList := []domain.BlockInterface{block1, block2, block3}
jsonBytes, err := blocks.SerializeBlocks(blockList)
if err != nil {
    log.Fatal(err)
}
```

## Deserializing Blocks

```go
import "github.com/vibezz/cms/internal/blocks"

// Deserialize a single block
var jsonBlock blocks.BlockJSON
if err := json.Unmarshal(data, &jsonBlock); err != nil {
    log.Fatal(err)
}
block, err := blocks.DeserializeBlock(jsonBlock)

// Deserialize multiple blocks
blockList, err := blocks.DeserializeBlocks(jsonData)
if err != nil {
    log.Fatal(err)
}
```

## Building a Block Tree

When blocks are stored flat in the database, you can rebuild the tree structure:

```go
// blocks is a flat list from the database
tree := blocks.BuildBlockTree(flatBlockList)

// Now tree contains only root blocks, each with their children
for _, rootBlock := range tree {
    fmt.Printf("Root: %s\n", rootBlock.BlockType())
    for _, child := range rootBlock.Children() {
        fmt.Printf("  Child: %s\n", child.BlockType())
    }
}
```

## Flattening a Block Tree

To store blocks in a flat database structure:

```go
// tree is a nested block structure
flatList := blocks.FlattenBlocks(tree)

// flatList now contains all blocks in depth-first order
for _, block := range flatList {
    // Save to database
    db.SaveBlock(block)
}
```

## Validating Blocks

```go
import "github.com/vibezz/cms/internal/blocks"

block := createSomeBlock()
if err := blocks.ValidateBlock(block); err != nil {
    fmt.Printf("Validation error: %v\n", err)
}
```

## Registering Custom Block Types

Plugins can register their own block types:

```go
import "github.com/vibezz/cms/internal/plugin"

func (p *MyPlugin) Register(deps *plugin.Deps) {
    // Register a custom block type
    deps.RegisterBlockType(plugin.BlockTypeDefinition{
        Slug:     "my_custom_block",
        Label:    "My Custom Block",
        Icon:     "star",
        Category: "custom",
        DefaultData: map[string]interface{}{
            "title": "",
            "items": []string{},
        },
        Fields: []plugin.BlockFieldDefinition{
            {Name: "title", Type: "text", Label: "Title", Required: true},
            {Name: "items", Type: "array", Label: "Items"},
        },
    })
    
    // Register an enricher for the block
    deps.RegisterBlockEnricher("my_custom_block", func(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {
        // Enrich the block data for public API
        data["enriched"] = true
        return data, nil
    })
}
```

## Block Styling

All blocks support styling:

```go
block.SetStyle(domain.BlockStyle{
    ClassName: "my-custom-class",
    InlineStyle: map[string]string{
        "color": "#333",
        "fontSize": "16px",
    },
    Responsive: &domain.ResponsiveBreakpoints{
        Mobile: &domain.BlockStyle{
            ClassName: "mobile-class",
            InlineStyle: map[string]string{
                "fontSize": "14px",
            },
        },
        Tablet: &domain.BlockStyle{
            ClassName: "tablet-class",
        },
        Desktop: &domain.BlockStyle{
            ClassName: "desktop-class",
            InlineStyle: map[string]string{
                "fontSize": "18px",
            },
        },
    },
})
```

## Container Configuration

Containers support various layout options:

```go
container := domain.NewContainerBlock()

// Layout types: "flex", "grid", "row", "column", "stack"
container.SetLayout("flex")

// Gap between children (CSS value)
container.SetGap("16px")
container.SetPadding("24px")

// Background (color, gradient, or image URL)
container.SetBackground("linear-gradient(135deg, #667eea 0%, #764ba2 100%)")

// Max width constraint
container.SetMaxWidth("1200px")
```

## Grid Configuration

Grids support responsive column counts:

```go
grid := domain.NewGridBlock()

// Default columns
grid.SetColumns(3)

// Gaps
grid.SetColumnGap("16px")
grid.SetRowGap("24px")

// Responsive columns
grid.SetResponsiveColumns(domain.ResponsiveColumns{
    Mobile:  1,  // 1 column on mobile (< 640px)
    Tablet:  2,  // 2 columns on tablet (640px - 1024px)
    Desktop: 3,  // 3 columns on desktop (> 1024px)
})
```

## JSON Format

Blocks are serialized to JSON in the following format:

```json
{
  "id": "uuid",
  "type": "container",
  "entryId": "uuid",
  "parentId": null,
  "sequence": 0,
  "data": {
    "layout": "flex",
    "gap": "16px",
    "padding": "24px"
  },
  "style": {
    "className": "my-class",
    "inlineStyle": {
      "backgroundColor": "#f5f5f5"
    },
    "responsive": {
      "mobile": {
        "className": "mobile-class"
      }
    }
  },
  "children": [
    {
      "id": "uuid",
      "type": "heading",
      "sequence": 0,
      "data": {
        "text": "Hello",
        "level": 1
      }
    }
  ]
}
```

## Database Schema

Blocks are stored in a flat structure in the database:

```sql
CREATE TABLE blocks (
    id UUID PRIMARY KEY,
    entry_id UUID NOT NULL REFERENCES entries(id),
    parent_id UUID REFERENCES blocks(id),
    block_type VARCHAR(50) NOT NULL,
    sequence INTEGER DEFAULT 0,
    block_data JSONB DEFAULT '{}',
    style JSONB DEFAULT '{}',
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    soft_deleted_at TIMESTAMP
);

CREATE INDEX idx_blocks_entry ON blocks(entry_id);
CREATE INDEX idx_blocks_parent ON blocks(parent_id);
```

To query blocks for an entry:

```go
query := domain.NewBlockQuery().
    SetEntryID(entryID).
    SetOrderBy("sequence").
    SetSortOrder("asc")

blocks, err := store.BlockList(ctx, query)
tree := blocks.BuildBlockTree(blocks)
```
