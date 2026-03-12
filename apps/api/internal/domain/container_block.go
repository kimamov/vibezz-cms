package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Ensure containerBlock implements ContainerBlockInterface
var _ ContainerBlockInterface = (*containerBlock)(nil)

// containerBlock is a block that can contain other blocks
type containerBlock struct {
	*block
}

// NewContainerBlock creates a new container block
func NewContainerBlock() ContainerBlockInterface {
	cb := &containerBlock{
		block: NewBlock().(*block),
	}
	cb.SetBlockType("container")
	cb.SetCreatedAt(time.Now().UTC())
	cb.SetUpdatedAt(time.Now().UTC())
	cb.SetSoftDeletedAt(SoftDeleteTime)

	// Set default layout config
	cb.SetLayout("flex")
	cb.SetGap("16px")
	cb.SetPadding("0px")
	cb.SetBackground("")
	cb.SetMaxWidth("")

	return cb
}

// Layout returns the layout type
func (cb *containerBlock) Layout() string {
	return cb.Get("layout")
}

// SetLayout sets the layout type (flex, grid, stack, row, column)
func (cb *containerBlock) SetLayout(layout string) {
	cb.Set("layout", layout)
}

// Gap returns the gap between child blocks
func (cb *containerBlock) Gap() string {
	return cb.Get("gap")
}

// SetGap sets the gap between child blocks
func (cb *containerBlock) SetGap(gap string) {
	cb.Set("gap", gap)
}

// Padding returns the padding
func (cb *containerBlock) Padding() string {
	return cb.Get("padding")
}

// SetPadding sets the padding
func (cb *containerBlock) SetPadding(padding string) {
	cb.Set("padding", padding)
}

// Background returns the background
func (cb *containerBlock) Background() string {
	return cb.Get("background")
}

// SetBackground sets the background (color, gradient, or image URL)
func (cb *containerBlock) SetBackground(background string) {
	cb.Set("background", background)
}

// MaxWidth returns the max width
func (cb *containerBlock) MaxWidth() string {
	return cb.Get("max_width")
}

// SetMaxWidth sets the max width
func (cb *containerBlock) SetMaxWidth(maxWidth string) {
	cb.Set("max_width", maxWidth)
}

// Clone creates a deep copy
func (cb *containerBlock) Clone() BlockInterface {
	cloned := NewContainerBlock()
	clonedBlock := cloned.(*containerBlock)

	clonedBlock.SetID(uuid.New())
	clonedBlock.SetEntryID(cb.EntryID())
	clonedBlock.SetSequence(cb.Sequence())
	if parentID := cb.ParentID(); parentID != nil {
		clonedBlock.SetParentID(parentID)
	}

	clonedBlock.SetLayout(cb.Layout())
	clonedBlock.SetGap(cb.Gap())
	clonedBlock.SetPadding(cb.Padding())
	clonedBlock.SetBackground(cb.Background())
	clonedBlock.SetMaxWidth(cb.MaxWidth())
	clonedBlock.SetStyle(cb.Style())
	clonedBlock.SetBlockData(cb.BlockData())

	// Clone children
	if len(cb.children) > 0 {
		clonedChildren := make([]BlockInterface, len(cb.children))
		for i, child := range cb.children {
			clonedChildren[i] = child.Clone()
		}
		clonedBlock.SetChildren(clonedChildren)
	}

	return clonedBlock
}

// ContainerConfig represents the configuration for a container block
type ContainerConfig struct {
	Layout     string `json:"layout"`
	Gap        string `json:"gap"`
	Padding    string `json:"padding"`
	Background string `json:"background"`
	MaxWidth   string `json:"maxWidth"`
}

// ToContainerConfig converts BlockData to ContainerConfig
func (cb *containerBlock) ToContainerConfig() ContainerConfig {
	return ContainerConfig{
		Layout:     cb.Layout(),
		Gap:        cb.Gap(),
		Padding:    cb.Padding(),
		Background: cb.Background(),
		MaxWidth:   cb.MaxWidth(),
	}
}

// FromContainerConfig sets the block data from ContainerConfig
func (cb *containerBlock) FromContainerConfig(config ContainerConfig) {
	cb.SetLayout(config.Layout)
	cb.SetGap(config.Gap)
	cb.SetPadding(config.Padding)
	cb.SetBackground(config.Background)
	cb.SetMaxWidth(config.MaxWidth)
}

// ValidateContainerConfig validates container configuration
func ValidateContainerConfig(config ContainerConfig) error {
	validLayouts := map[string]bool{"flex": true, "grid": true, "stack": true, "row": true, "column": true}
	if !validLayouts[config.Layout] {
		return fmt.Errorf("invalid layout: %s", config.Layout)
	}
	return nil
}
