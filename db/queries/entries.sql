-- name: CreateEntry :one
INSERT INTO entries (id, content_type_id, title, slug, path, parent_id, author_id, status, fields, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING id, content_type_id, title, slug, path, parent_id, author_id, status, fields, published_at, created_at, updated_at;

-- name: GetEntryByID :one
SELECT id, content_type_id, title, slug, path, parent_id, author_id, status, fields, published_at, created_at, updated_at
FROM entries
WHERE id = $1;

-- name: GetEntryByPath :one
SELECT id, content_type_id, title, slug, path, parent_id, author_id, status, fields, published_at, created_at, updated_at
FROM entries
WHERE path = $1 AND status = 'published';

-- name: GetEntryPath :one
SELECT path FROM entries WHERE id = $1;

-- name: ListEntries :many
SELECT id, content_type_id, title, slug, path, parent_id, author_id, status, fields, published_at, created_at, updated_at
FROM entries
ORDER BY created_at DESC;

-- name: ListEntriesByContentType :many
SELECT id, content_type_id, title, slug, path, parent_id, author_id, status, fields, published_at, created_at, updated_at
FROM entries
WHERE content_type_id = $1
ORDER BY created_at DESC;

-- name: ListEntriesByStatus :many
SELECT id, content_type_id, title, slug, path, parent_id, author_id, status, fields, published_at, created_at, updated_at
FROM entries
WHERE status = $1
ORDER BY created_at DESC;

-- name: UpdateEntry :exec
UPDATE entries
SET title = $2, slug = $3, path = $4, fields = $5, updated_at = $6
WHERE id = $1;

-- name: PublishEntry :exec
UPDATE entries
SET status = 'published', published_at = $2, updated_at = $3
WHERE id = $1;

-- name: UnpublishEntry :exec
UPDATE entries
SET status = 'draft', published_at = NULL, updated_at = $2
WHERE id = $1;

-- name: DeleteEntry :exec
DELETE FROM entries WHERE id = $1;

-- name: GetNavigationItems :many
SELECT id, title, slug, path, parent_id
FROM entries
WHERE status = 'published'
ORDER BY path;
