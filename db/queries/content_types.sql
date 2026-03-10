-- name: ListContentTypes :many
SELECT id, name, slug, fields, created_at, updated_at
FROM content_types
ORDER BY name;

-- name: GetContentTypeByID :one
SELECT id, name, slug, fields, created_at, updated_at
FROM content_types
WHERE id = $1;

-- name: GetContentTypeBySlug :one
SELECT id, name, slug, fields, created_at, updated_at
FROM content_types
WHERE slug = $1;

-- name: CreateContentType :one
INSERT INTO content_types (id, name, slug, fields, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, name, slug, fields, created_at, updated_at;

-- name: UpdateContentType :exec
UPDATE content_types
SET name = $2, fields = $3, updated_at = $4
WHERE id = $1;

-- name: DeleteContentType :exec
DELETE FROM content_types WHERE id = $1;
