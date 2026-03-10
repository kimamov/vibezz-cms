-- name: CreateMediaFile :one
INSERT INTO media_files (id, filename, mime_type, size, storage_path, uploader_id, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, filename, mime_type, size, storage_path, uploader_id, created_at;

-- name: GetMediaFileByID :one
SELECT id, filename, mime_type, size, storage_path, uploader_id, created_at
FROM media_files
WHERE id = $1;

-- name: ListMediaFiles :many
SELECT id, filename, mime_type, size, storage_path, uploader_id, created_at
FROM media_files
ORDER BY created_at DESC;

-- name: DeleteMediaFile :exec
DELETE FROM media_files WHERE id = $1;
