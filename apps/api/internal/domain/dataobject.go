package domain

// DataObject provides dirty tracking capabilities for entities
// Following the cmsstore pattern
type DataObject struct {
	data        map[string]string
	dataChanged map[string]string
	isDirty     bool
}

// NewDataObject creates a new DataObject instance
func NewDataObject() *DataObject {
	return &DataObject{
		data:        make(map[string]string),
		dataChanged: make(map[string]string),
		isDirty:     false,
	}
}

// Hydrate populates the data object with existing data
func (o *DataObject) Hydrate(data map[string]string) {
	o.data = data
	o.dataChanged = make(map[string]string)
	o.isDirty = false
}

// Data returns all data
func (o *DataObject) Data() map[string]string {
	return o.data
}

// DataChanged returns only the changed fields
func (o *DataObject) DataChanged() map[string]string {
	return o.dataChanged
}

// MarkAsNotDirty clears the dirty state
func (o *DataObject) MarkAsNotDirty() {
	o.dataChanged = make(map[string]string)
	o.isDirty = false
}

// Get retrieves a value by key
func (o *DataObject) Get(key string) string {
	return o.data[key]
}

// Set sets a value and marks the field as changed
func (o *DataObject) Set(key string, value string) {
	oldValue := o.data[key]
	if oldValue != value {
		o.data[key] = value
		o.dataChanged[key] = value
		o.isDirty = true
	}
}

// IsDirty returns true if any field has been modified
func (o *DataObject) IsDirty() bool {
	return o.isDirty
}

// SetMultiple sets multiple values at once
func (o *DataObject) SetMultiple(data map[string]string) {
	for key, value := range data {
		o.Set(key, value)
	}
}
