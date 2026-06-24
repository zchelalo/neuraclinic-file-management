CREATE TABLE files (
  id uuid PRIMARY KEY,
  original_name varchar(255) NOT NULL,
  storage_path varchar(1024) NOT NULL,
  mime_type varchar(100) NOT NULL,
  size_bytes bigint NOT NULL,
  checksum varchar(128) NOT NULL DEFAULT '',
  service_origin varchar(50) NOT NULL,
  is_public boolean NOT NULL DEFAULT false,
  status varchar(50) NOT NULL,
  storage_provider varchar(50) NOT NULL DEFAULT 's3',
  uploaded_by uuid NOT NULL,
  uploaded_at timestamptz,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),
  deleted_at timestamptz
);

CREATE INDEX idx_files_uploaded_by_created
  ON files (uploaded_by, created_at DESC, id DESC)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_files_service_origin_created
  ON files (service_origin, created_at DESC, id DESC)
  WHERE deleted_at IS NULL;

CREATE INDEX idx_files_status_created
  ON files (status, created_at DESC)
  WHERE deleted_at IS NULL;

