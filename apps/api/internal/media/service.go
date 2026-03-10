package media

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vibezz/cms/internal/content"
)

type Service struct {
	pool      *pgxpool.Pool
	uploadDir string
}

func NewService(pool *pgxpool.Pool, uploadDir string) *Service {
	return &Service{pool: pool, uploadDir: uploadDir}
}

func (s *Service) DiskPath(storagePath string) string {
	return filepath.Join(s.uploadDir, storagePath)
}

func (s *Service) Upload(ctx context.Context, header *multipart.FileHeader, uploaderID uuid.UUID) (*content.MediaFile, error) {
	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	id := uuid.New()
	ext := filepath.Ext(header.Filename)
	storageName := id.String() + ext
	storagePath := storageName

	src, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dstPath := filepath.Join(s.uploadDir, storageName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	buf := make([]byte, 32*1024)
	var written int64
	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			nw, writeErr := dst.Write(buf[:n])
			written += int64(nw)
			if writeErr != nil {
				return nil, writeErr
			}
		}
		if readErr != nil {
			break
		}
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	now := time.Now()
	_, err = s.pool.Exec(ctx,
		`INSERT INTO media_files (id, filename, mime_type, size, storage_path, uploader_id, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		id, header.Filename, mimeType, written, storagePath, uploaderID, now)
	if err != nil {
		os.Remove(dstPath)
		return nil, err
	}

	return &content.MediaFile{
		ID: id, Filename: header.Filename,
		MimeType: mimeType, Size: written,
		StoragePath: storagePath, UploaderID: uploaderID,
		CreatedAt: now,
	}, nil
}

func (s *Service) List(ctx context.Context) ([]content.MediaFile, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT id, filename, mime_type, size, storage_path, uploader_id, created_at
		 FROM media_files ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []content.MediaFile
	for rows.Next() {
		var f content.MediaFile
		if err := rows.Scan(&f.ID, &f.Filename, &f.MimeType, &f.Size, &f.StoragePath, &f.UploaderID, &f.CreatedAt); err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	if files == nil {
		files = []content.MediaFile{}
	}
	return files, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*content.MediaFile, error) {
	row := s.pool.QueryRow(ctx,
		`SELECT id, filename, mime_type, size, storage_path, uploader_id, created_at
		 FROM media_files WHERE id = $1`, id)

	var f content.MediaFile
	if err := row.Scan(&f.ID, &f.Filename, &f.MimeType, &f.Size, &f.StoragePath, &f.UploaderID, &f.CreatedAt); err != nil {
		return nil, err
	}
	return &f, nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	file, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	diskPath := s.DiskPath(file.StoragePath)
	os.Remove(diskPath)

	_, err = s.pool.Exec(ctx, `DELETE FROM media_files WHERE id = $1`, id)
	return err
}
