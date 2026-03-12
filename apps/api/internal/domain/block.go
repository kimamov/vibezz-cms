package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Ensure block implements BlockInterface
var _ BlockInterface = (*block)(nil)

// block is the concrete implementation of BlockInterface
type block struct {
	*DataObject
	children []BlockInterface
}

// NewBlock creates a new Block instance with default values
func NewBlock() BlockInterface {
	b := &block{
		DataObject: NewDataObject(),
		children:   []BlockInterface{},
	}
	b.SetID(uuid.New())
	b.SetCreatedAt(time.Now().UTC())
	b.SetUpdatedAt(time.Now().UTC())
	b.SetSoftDeletedAt(SoftDeleteTime)
	b.SetBlockData(map[string]interface{}{})
	b.SetStyle(BlockStyle{})
	return b
}

// NewBlockFromData creates a Block from existing data
func NewBlockFromData(data map[string]string) BlockInterface {
	b := &block{
		DataObject: NewDataObject(),
		children:   []BlockInterface{},
	}
	b.Hydrate(data)
	return b
}

// ID returns the block ID
func (b *block) ID() uuid.UUID {
	return uuid.MustParse(b.Get(ColumnID))
}

// SetID sets the block ID
func (b *block) SetID(id uuid.UUID) {
	b.Set(ColumnID, id.String())
}

// CreatedAt returns the creation time
func (b *block) CreatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, b.Get(ColumnCreatedAt))
	return t
}

// SetCreatedAt sets the creation time
func (b *block) SetCreatedAt(createdAt time.Time) {
	b.Set(ColumnCreatedAt, createdAt.Format(time.RFC3339))
}

// UpdatedAt returns the last update time
func (b *block) UpdatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, b.Get(ColumnUpdatedAt))
	return t
}

// SetUpdatedAt sets the last update time
func (b *block) SetUpdatedAt(updatedAt time.Time) {
	b.Set(ColumnUpdatedAt, updatedAt.Format(time.RFC3339))
}

// SoftDeletedAt returns the soft delete timestamp
func (b *block) SoftDeletedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, b.Get(ColumnSoftDeletedAt))
	return t
}

// SetSoftDeletedAt sets the soft delete timestamp
func (b *block) SetSoftDeletedAt(softDeletedAt time.Time) {
	b.Set(ColumnSoftDeletedAt, softDeletedAt.Format(time.RFC3339))
}

// IsSoftDeleted returns true if the block is soft deleted
func (b *block) IsSoftDeleted() bool {
	return IsSoftDeleted(b.SoftDeletedAt())
}

// BlockType returns the block type
func (b *block) BlockType() string {
	return b.Get("block_type")
}

// SetBlockType sets the block type
func (b *block) SetBlockType(blockType string) {
	b.Set("block_type", blockType)
}

// EntryID returns the entry ID
func (b *block) EntryID() uuid.UUID {
	return uuid.MustParse(b.Get("entry_id"))
}

// SetEntryID sets the entry ID
func (b *block) SetEntryID(entryID uuid.UUID) {
	b.Set("entry_id", entryID.String())
}

// ParentID returns the parent block ID
func (b *block) ParentID() *uuid.UUID {
	parentIDStr := b.Get("parent_id")
	if parentIDStr == "" {
		return nil
	}
	id := uuid.MustParse(parentIDStr)
	return &id
}

// SetParentID sets the parent block ID
func (b *block) SetParentID(parentID *uuid.UUID) {
	if parentID == nil {
		b.Set("parent_id", "")
	} else {
		b.Set("parent_id", parentID.String())
	}
}

// Sequence returns the block sequence/order
func (b *block) Sequence() int {
	seq := b.Get("sequence")
	if seq == "" {
		return 0
	}
	var s int
	fmt.Sscanf(seq, "%d", &s)
	return s
}

// SetSequence sets the block sequence
func (b *block) SetSequence(sequence int) {
	b.Set("sequence", fmt.Sprintf("%d", sequence))
}

// BlockData returns the block-specific data
func (b *block) BlockData() map[string]interface{} {
	dataStr := b.Get("block_data")
	if dataStr == "" {
		return map[string]interface{}{}
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(dataStr), &data); err != nil {
		return map[string]interface{}{}
	}
	return data
}

// SetBlockData sets the block-specific data
func (b *block) SetBlockData(data map[string]interface{}) error {
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}
	b.Set("block_data", string(dataJSON))
	return nil
}

// Children returns child blocks
func (b *block) Children() []BlockInterface {
	return b.children
}

// SetChildren sets child blocks
func (b *block) SetChildren(children []BlockInterface) {
	b.children = children
}

// Style returns the block style
func (b *block) Style() BlockStyle {
	styleStr := b.Get("style")
	if styleStr == "" {
		return BlockStyle{}
	}

	var style BlockStyle
	if err := json.Unmarshal([]byte(styleStr), &style); err != nil {
		return BlockStyle{}
	}
	return style
}

// SetStyle sets the block style
func (b *block) SetStyle(style BlockStyle) {
	styleJSON, _ := json.Marshal(style)
	b.Set("style", string(styleJSON))
}

// IsContainer returns true if this block can contain other blocks
func (b *block) IsContainer() bool {
	blockType := b.BlockType()
	return blockType == "container" || blockType == "grid" || blockType == "section"
}

// Clone creates a deep copy of the block
func (b *block) Clone() BlockInterface {
	cloned := NewBlock()
	cloned.SetBlockType(b.BlockType())
	cloned.SetEntryID(b.EntryID())
	cloned.SetSequence(b.Sequence())

	if parentID := b.ParentID(); parentID != nil {
		cloned.SetParentID(parentID)
	}

	data := b.BlockData()
	cloned.SetBlockData(data)
	cloned.SetStyle(b.Style())

	// Clone children
	if len(b.children) > 0 {
		clonedChildren := make([]BlockInterface, len(b.children))
		for i, child := range b.children {
			clonedChildren[i] = child.Clone()
		}
		cloned.SetChildren(clonedChildren)
	}

	return cloned
}
