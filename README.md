# Neuraclinic File Management Microservice

Go gRPC file-management service for Neuraclinic.

It stores file metadata in PostgreSQL, creates pre-signed S3 URLs, and publishes file lifecycle events to RabbitMQ. Local development uses MinIO as the S3-compatible storage backend.

## Local Setup

Run from `neuraclinic-file-management`:

```bash
make create-envs
make tls-generate-dev
```

Run shared services from the root `neuraclinic` repository first:

```bash
cd ../neuraclinic
make compose-detached
cd ../neuraclinic-file-management
make create-network
make compose-build
```

The local compose stack starts:

- `neuraclinic-file-management` on host port `8004`
- PostgreSQL on host port `5437`
- Adminer on `http://localhost:8087`
- MinIO S3 API on `http://localhost:9000`
- MinIO console on `http://localhost:9001`

`neuraclinic-rabbitmq` must also be running on `neuraclinic-network` because this service publishes `FileStatusChangedEvent` there.

Relevant env vars:

- `RABBITMQ_URL=amqp://guest:guest@neuraclinic-rabbitmq:5672/`
- `RABBITMQ_EXCHANGE=neuraclinic.events`

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
