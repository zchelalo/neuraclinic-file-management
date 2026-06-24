# Neuraclinic File Management Microservice

Go gRPC file-management service for Neuraclinic.

It stores file metadata in PostgreSQL and creates pre-signed S3 URLs. Local development uses MinIO as the S3-compatible storage backend.

## Local Setup

Run from `neuraclinic-file-management`:

```bash
make create-envs
make tls-generate-dev
make compose-build
```

The local compose stack starts:

- `neuraclinic-file-management` on host port `8004`
- PostgreSQL on host port `5437`
- Adminer on `http://localhost:8087`
- MinIO S3 API on `http://localhost:9000`
- MinIO console on `http://localhost:9001`

## Useful Commands

```bash
make proto
make sqlc
make test
make build
make migrate-up
make compose
make compose-down
```

## Storage

Local MinIO credentials:

- access key: `neuraclinic`
- secret key: `neuraclinic-secret`
- bucket: `neuraclinic-local-files`

Terraform for an AWS S3 bucket lives in `infra/aws-s3`.

