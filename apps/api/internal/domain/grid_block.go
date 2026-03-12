package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Ensure gridBlock implements GridBlockInterface
var _ GridBlockInterface = (*gridBlock)(nil)

// gridBlock is a block that arranges children in a responsive grid
type gridBlock struct {
	*block
}

// NewGridBlock creates a new grid block
func NewGridBlock() GridBlockInterface {
	gb := &gridBlock{
		block: NewBlock().(*block),
	}
	gb.SetBlockType("grid")
	gb.SetCreatedAt(time.Now().UTC())
	gb.SetUpdatedAt(time.Now().UTC())
	gb.SetSoftDeletedAt(SoftDeleteTime)

	// Set default grid config
	gb.SetColumns(3)
	gb.SetColumnGap("16px")
	gb.SetRowGap("16px")
	gb.SetResponsiveColumns(ResponsiveColumns{
		Mobile:  1,
		Tablet:  2,
		Desktop: 3,
	})

	return gb
}

// Columns returns the number of columns
func (gb *gridBlock) Columns() int {
	cols := gb.Get("columns")
	if cols == "" {
		return 3
	}
	var c int
	fmt.Sscanf(cols, "%d", &c)
	return c
}

// SetColumns sets the number of columns
func (gb *gridBlock) SetColumns(columns int) {
	if columns < 1 {
		columns = 1
	}
	if columns > 12 {
		columns = 12
	}
	gb.Set("columns", fmt.Sprintf("%d", columns))
}

// ColumnGap returns the gap between columns
func (gb *gridBlock) ColumnGap() string {
	return gb.Get("column_gap")
}

// SetColumnGap sets the gap between columns
func (gb *gridBlock) SetColumnGap(gap string) {
	gb.Set("column_gap", gap)
}

// RowGap returns the gap between rows
func (gb *gridBlock) RowGap() string {
	return gb.Get("row_gap")
}

// SetRowGap sets the gap between rows
func (gb *gridBlock) SetRowGap(gap string) {
	gb.Set("row_gap", gap)
}

// ResponsiveColumns returns the responsive column configuration
func (gb *gridBlock) ResponsiveColumns() ResponsiveColumns {
	colsStr := gb.Get("responsive_columns")
	if colsStr == "" {
		return ResponsiveColumns{Mobile: 1, Tablet: 2, Desktop: 3}
	}

	var cols ResponsiveColumns
	fmt.Sscanf(colsStr, "%d,%d,%d", &cols.Mobile, &cols.Tablet, &cols.Desktop)
	return cols
}

// SetResponsiveColumns sets the responsive column configuration
func (gb *gridBlock) SetResponsiveColumns(cols ResponsiveColumns) {
	// Validate
	if cols.Mobile < 1 {
		cols.Mobile = 1
	}
	if cols.Tablet < 1 {
		cols.Tablet = 1
	}
	if cols.Desktop < 1 {
		cols.Desktop = 1
	}
	if cols.Mobile > 12 {
		cols.Mobile = 12
	}
	if cols.Tablet > 12 {
		cols.Tablet = 12
	}
	if cols.Desktop > 12 {
		cols.Desktop = 12
	}

	gb.Set("responsive_columns", fmt.Sprintf("%d,%d,%d", cols.Mobile, cols.Tablet, cols.Desktop))
}

// Clone creates a deep copy
func (gb *gridBlock) Clone() BlockInterface {
	cloned := NewGridBlock()
	clonedBlock := cloned.(*gridBlock)

	clonedBlock.SetID(uuid.New())
	clonedBlock.SetEntryID(gb.EntryID())
	clonedBlock.SetSequence(gb.Sequence())
	if parentID := gb.ParentID(); parentID != nil {
		clonedBlock.SetParentID(parentID)
	}

	clonedBlock.SetColumns(gb.Columns())
	clonedBlock.SetColumnGap(gb.ColumnGap())
	clonedBlock.SetRowGap(gb.RowGap())
	clonedBlock.SetResponsiveColumns(gb.ResponsiveColumns())
	clonedBlock.SetStyle(gb.Style())
	clonedBlock.SetBlockData(gb.BlockData())

	// Clone children
	if len(gb.children) > 0 {
		clonedChildren := make([]BlockInterface, len(gb.children))
		for i, child := range gb.children {
			clonedChildren[i] = child.Clone()
		}
		clonedBlock.SetChildren(clonedChildren)
	}

	return clonedBlock
}

// GridConfig represents the configuration for a grid block
type GridConfig struct {
	Columns           int               `json:"columns"`
	ColumnGap         string            `json:"columnGap"`
	RowGap            string            `json:"rowGap"`
	ResponsiveColumns ResponsiveColumns `json:"responsiveColumns"`
}

// ToGridConfig converts BlockData to GridConfig
func (gb *gridBlock) ToGridConfig() GridConfig {
	return GridConfig{
		Columns:           gb.Columns(),
		ColumnGap:         gb.ColumnGap(),
		RowGap:            gb.RowGap(),
		ResponsiveColumns: gb.ResponsiveColumns(),
	}
}

// FromGridConfig sets the block data from GridConfig
func (gb *gridBlock) FromGridConfig(config GridConfig) {
	gb.SetColumns(config.Columns)
	gb.SetColumnGap(config.ColumnGap)
	gb.SetRowGap(config.RowGap)
	gb.SetResponsiveColumns(config.ResponsiveColumns)
}

// ValidateGridConfig validates grid configuration
func ValidateGridConfig(config GridConfig) error {
	if config.Columns < 1 || config.Columns > 12 {
		return fmt.Errorf("columns must be between 1 and 12, got %d", config.Columns)
	}
	if config.ResponsiveColumns.Mobile < 1 || config.ResponsiveColumns.Mobile > 12 {
		return fmt.Errorf("mobile columns must be between 1 and 12")
	}
	if config.ResponsiveColumns.Tablet < 1 || config.ResponsiveColumns.Tablet > 12 {
		return fmt.Errorf("tablet columns must be between 1 and 12")
	}
	if config.ResponsiveColumns.Desktop < 1 || config.ResponsiveColumns.Desktop > 12 {
		return fmt.Errorf("desktop columns must be between 1 and 12")
	}
	return nil
}
