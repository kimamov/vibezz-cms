package domain

import "errors"

// query is the base implementation of QueryInterface
type query struct {
	properties map[string]interface{}
}

// NewQuery creates a new Query instance
func NewQuery() QueryInterface {
	return &query{
		properties: make(map[string]interface{}),
	}
}

// Validate validates the query
func (q *query) Validate() error {
	if q.HasLimit() && q.Limit() < 0 {
		return errors.New("limit cannot be negative")
	}
	if q.HasOffset() && q.Offset() < 0 {
		return errors.New("offset cannot be negative")
	}
	return nil
}

// HasLimit returns true if limit is set
func (q *query) HasLimit() bool {
	_, ok := q.properties["limit"]
	return ok
}

// Limit returns the limit
func (q *query) Limit() int {
	if val, ok := q.properties["limit"].(int); ok {
		return val
	}
	return 0
}

// SetLimit sets the limit
func (q *query) SetLimit(limit int) QueryInterface {
	q.properties["limit"] = limit
	return q
}

// HasOffset returns true if offset is set
func (q *query) HasOffset() bool {
	_, ok := q.properties["offset"]
	return ok
}

// Offset returns the offset
func (q *query) Offset() int {
	if val, ok := q.properties["offset"].(int); ok {
		return val
	}
	return 0
}

// SetOffset sets the offset
func (q *query) SetOffset(offset int) QueryInterface {
	q.properties["offset"] = offset
	return q
}

// HasOrderBy returns true if order by is set
func (q *query) HasOrderBy() bool {
	_, ok := q.properties["order_by"]
	return ok
}

// OrderBy returns the order by field
func (q *query) OrderBy() string {
	if val, ok := q.properties["order_by"].(string); ok {
		return val
	}
	return ""
}

// SetOrderBy sets the order by field
func (q *query) SetOrderBy(orderBy string) QueryInterface {
	q.properties["order_by"] = orderBy
	return q
}

// HasSortOrder returns true if sort order is set
func (q *query) HasSortOrder() bool {
	_, ok := q.properties["sort_order"]
	return ok
}

// SortOrder returns the sort order
func (q *query) SortOrder() string {
	if val, ok := q.properties["sort_order"].(string); ok {
		return val
	}
	return "asc"
}

// SetSortOrder sets the sort order
func (q *query) SetSortOrder(sortOrder string) QueryInterface {
	q.properties["sort_order"] = sortOrder
	return q
}
