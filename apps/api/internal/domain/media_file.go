package domain

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Ensure mediaFile implements MediaFileInterface
var _ MediaFileInterface = (*mediaFile)(nil)

// mediaFile is the concrete implementation of MediaFileInterface
type mediaFile struct {
	*DataObject
}

// NewMediaFile creates a new MediaFile instance with default values
func NewMediaFile() MediaFileInterface {
	m := &mediaFile{
		DataObject: NewDataObject(),
	}
	m.SetID(uuid.New())
	m.SetCreatedAt(time.Now().UTC())
	m.SetUpdatedAt(time.Now().UTC())
	m.SetSoftDeletedAt(SoftDeleteTime)
	return m
}

// NewMediaFileFromData creates a MediaFile from existing data
func NewMediaFileFromData(data map[string]string) MediaFileInterface {
	m := &mediaFile{
		DataObject: NewDataObject(),
	}
	m.Hydrate(data)
	return m
}

// ID returns the media file ID
func (m *mediaFile) ID() uuid.UUID {
	return uuid.MustParse(m.Get(ColumnID))
}

// SetID sets the media file ID
func (m *mediaFile) SetID(id uuid.UUID) {
	m.Set(ColumnID, id.String())
}

// CreatedAt returns the creation time
func (m *mediaFile) CreatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, m.Get(ColumnCreatedAt))
	return t
}

// SetCreatedAt sets the creation time
func (m *mediaFile) SetCreatedAt(createdAt time.Time) {
	m.Set(ColumnCreatedAt, createdAt.Format(time.RFC3339))
}

// UpdatedAt returns the last update time
func (m *mediaFile) UpdatedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, m.Get(ColumnUpdatedAt))
	return t
}

// SetUpdatedAt sets the last update time
func (m *mediaFile) SetUpdatedAt(updatedAt time.Time) {
	m.Set(ColumnUpdatedAt, updatedAt.Format(time.RFC3339))
}

// SoftDeletedAt returns the soft delete timestamp
func (m *mediaFile) SoftDeletedAt() time.Time {
	t, _ := time.Parse(time.RFC3339, m.Get(ColumnSoftDeletedAt))
	return t
}

// SetSoftDeletedAt sets the soft delete timestamp
func (m *mediaFile) SetSoftDeletedAt(softDeletedAt time.Time) {
	m.Set(ColumnSoftDeletedAt, softDeletedAt.Format(time.RFC3339))
}

// IsSoftDeleted returns true if the media file is soft deleted
func (m *mediaFile) IsSoftDeleted() bool {
	return IsSoftDeleted(m.SoftDeletedAt())
}

// Filename returns the filename
func (m *mediaFile) Filename() string {
	return m.Get("filename")
}

// SetFilename sets the filename
func (m *mediaFile) SetFilename(filename string) {
	m.Set("filename", filename)
}

// MimeType returns the MIME type
func (m *mediaFile) MimeType() string {
	return m.Get("mime_type")
}

// SetMimeType sets the MIME type
func (m *mediaFile) SetMimeType(mimeType string) {
	m.Set("mime_type", mimeType)
}

// Size returns the file size
func (m *mediaFile) Size() int64 {
	// Parse int64 from string storage
	sizeStr := m.Get("size")
	if sizeStr == "" {
		return 0
	}
	var size int64
	fmt.Sscanf(sizeStr, "%d", &size)
	return size
}

// SetSize sets the file size
func (m *mediaFile) SetSize(size int64) {
	m.Set("size", fmt.Sprintf("%d", size))
}

// StoragePath returns the storage path
func (m *mediaFile) StoragePath() string {
	return m.Get("storage_path")
}

// SetStoragePath sets the storage path
func (m *mediaFile) SetStoragePath(path string) {
	m.Set("storage_path", path)
}

// UploaderID returns the uploader ID
func (m *mediaFile) UploaderID() uuid.UUID {
	return uuid.MustParse(m.Get("uploader_id"))
}

// SetUploaderID sets the uploader ID
func (m *mediaFile) SetUploaderID(uploaderID uuid.UUID) {
	m.Set("uploader_id", uploaderID.String())
}
