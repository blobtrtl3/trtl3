# TRTL3

<img align="right" width="180px" src="https://github.com/blobtrtl3/docs/blob/main/logo/logo.svg" alt="trtl3 logo">

[![publish status](https://github.com/blobtrtl3/trtl3/actions/workflows/publish.yml/badge.svg?branch=main)](https://github.com/blobtrtl3/trtl3/actions/workflows/publish.yml)
[![dockerhub](https://img.shields.io/docker/pulls/nothiaki/trtl3)](https://hub.docker.com/r/nothiaki/trtl3)
[![go report](https://goreportcard.com/badge/github.com/blobtrtl3/trtl3)](https://goreportcard.com/report/github.com/blobtrtl3/trtl3)
[![APGL-3.0](https://img.shields.io/badge/license-AGPL--3.0-blue.svg)](LICENSE)
[![PRs welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![download](https://img.shields.io/badge/download-latest-brightgreen)](https://hub.docker.com/r/nothiaki/trtl3)

**TRTL3**🐢 is a file storage service written in [GO](https://go.dev/) designed for developers who want to build things without
the complexity of cloud services or external dependencies.
If you're building something with files uploads, downloads, and blobs organization — TRTL3 is a great starting point!

---

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Features](#features)

---

## Getting Started

After Install Docker/Docker Compose you can install TRTL3 using this command:

```bash
docker run -d --name trtl3 -p 7713:7713 -e TOKEN=your_secret_token nothiaki/trtl3:latest
```

So the service will be running on `http://localhost:7713/` and you can use by REST or use SDK's.

If you want to use docker compose you can add it on the file:

```yaml
services:

  trtl3:
    image: nothiaki/trtl3:latest
    ports:
      - 7713:7713
    environment:
      - TOKEN=your_token_here
      - WORKERS=10
```

### 🔧 Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TOKEN` | Authentication token for API access | `trtl3` |
| `WORKERS` | Number of worker threads for file processing | `10` |
| `JOB_INTERVAL` | Cleanup job interval in minutes | `5` |

---

## Usage

You can use TRTL3 by REST api or our official SDKs:
Go on our [website](https://trtl3.store) to learn more.

SDKs:

- [GO](https://github.com/blobtrtl3/trtl3-go)
... soon in nodejs and java

### 📝 API Examples

#### Upload a File
```bash
curl -X POST http://localhost:7713/blobs \
  -H "Authorization: Bearer your_token_here" \
  -F "bucket=my-bucket" \
  -F "blob=@/path/to/your/file.jpg"
```

#### List Files in Bucket
```bash
curl -X GET "http://localhost:7713/blobs?bucket=my-bucket" \
  -H "Authorization: Bearer your_token_here"
```

#### Create Signed URL
```bash
curl -X POST http://localhost:7713/blobs/sign \
  -H "Authorization: Bearer your_token_here" \
  -H "Content-Type: application/json" \
  -d '{
    "bucket": "my-bucket",
    "id": "file-id",
    "ttl": 60,
    "once": false
  }'
```

#### Access File via Signed URL
```bash
curl -X GET "http://localhost:7713/b?sign=your_signed_token"
```

#### Download File Directly
```bash
curl -X GET "http://localhost:7713/blobs/download/my-bucket/file-id" \
  -H "Authorization: Bearer your_token_here" \
  -o downloaded_file.jpg
```

#### Health Check
```bash
curl -X GET "http://localhost:7713/health"
```

---

## Features

### 🚀 Core Features

- **File Upload & Storage**: Upload files with automatic ID generation and metadata tracking
- **Bucket Organization**: Organize files into logical buckets (like folders)
- **Signed URLs**: Generate temporary, secure URLs for file access without authentication
- **REST API**: Complete RESTful API for all operations
- **Authentication**: Token-based authentication for secure access
- **Automatic Cleanup**: Background jobs to clean orphaned files and expired signatures

### 🔧 Technical Features

- **Async Processing**: Queue-based file processing with configurable workers
- **Embedded Database**: Uses DuckDB for metadata storage (no external dependencies)
- **Docker Ready**: Easy deployment with Docker and Docker Compose
- **File System Storage**: Simple file-based storage with automatic directory creation
- **TTL Support**: Configurable expiration for signed URLs (1-1440 minutes)
- **One-time Access**: Optional single-use URLs for sensitive files

### 📡 API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/blobs` | Upload a file |
| `GET` | `/blobs` | List files by bucket |
| `GET` | `/blobs/:bucket/:id` | Get file metadata |
| `DELETE` | `/blobs/:bucket/:id` | Delete a file |
| `GET` | `/blobs/download/:bucket/:id` | Download file directly |
| `POST` | `/blobs/sign` | Create signed URL |
| `GET` | `/b?sign=TOKEN` | Access file via signed URL |
| `GET` | `/health` | Health check endpoint |

### 🛡️ Security Features

- **Token Authentication**: Bearer token authentication for API access
- **Signed URLs**: Time-limited, secure access to files
- **Input Validation**: File type and size validation
- **Secure File Naming**: Automatic secure file naming to prevent conflicts
