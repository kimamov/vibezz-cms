package domain

import "github.com/google/uuid"

// blockQuery is the implementation of BlockQueryInterface
type blockQuery struct {
	*query
}

// NewBlockQuery creates a new BlockQuery instance
func NewBlockQuery() BlockQueryInterface {
	return &blockQuery{
		query: NewQuery().(*query),
	}
}

// HasEntryID returns true if entry ID is set
func (q *blockQuery) HasEntryID() bool {
	_, ok := q.properties["entry_id"]
	return ok
}

// EntryID returns the entry ID
func (q *blockQuery) EntryID() uuid.UUID {
	if val, ok := q.properties["entry_id"].(uuid.UUID); ok {
		return val
	}
	return uuid.Nil
}

// SetEntryID sets the entry ID
func (q *blockQuery) SetEntryID(entryID uuid.UUID) BlockQueryInterface {
	q.properties["entry_id"] = entryID
	return q
}

// HasParentID returns true if parent ID is set
func (q *blockQuery) HasParentID() bool {
	_, ok := q.properties["parent_id"]
	return ok
}

// ParentID returns the parent ID
func (q *blockQuery) ParentID() *uuid.UUID {
	if val, ok := q.properties["parent_id"].(*uuid.UUID); ok {
		return val
	}
	return nil
}

// SetParentID sets the parent ID
func (q *blockQuery) SetParentID(parentID *uuid.UUID) BlockQueryInterface {
	q.properties["parent_id"] = parentID
	return q
}

// HasBlockType returns true if block type is set
func (q *blockQuery) HasBlockType() bool {
	_, ok := q.properties["block_type"]
	return ok
}

// BlockType returns the block type
func (q *blockQuery) BlockType() string {
	if val, ok := q.properties["block_type"].(string); ok {
		return val
	}
	return ""
}

// SetBlockType sets the block type
func (q *blockQuery) SetBlockType(blockType string) BlockQueryInterface {
	q.properties["block_type"] = blockType
	return q
}
