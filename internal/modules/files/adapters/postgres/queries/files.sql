-- name: CreateFile :one
INSERT INTO files (
  id,
  original_name,
  storage_path,
  mime_type,
  size_bytes,
  checksum,
  service_origin,
  is_public,
  status,
  storage_provider,
  uploaded_by,
  created_at,
  updated_at
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $12
)
RETURNING *;

-- name: GetFileByID :one
SELECT *
FROM files
WHERE id = $1
  AND deleted_at IS NULL;

-- name: ConfirmFileUpload :one
UPDATE files
SET
  status = $2,
  uploaded_at = CASE
    WHEN $2 = 'FILE_STATUS_AVAILABLE' THEN $3
    ELSE uploaded_at
  END,
  updated_at = $3
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

