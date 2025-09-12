# TRTL3

<img align="right" width="180px" src="https://github.com/blobtrtl3/docs/blob/main/logo/logo.svg" alt="trtl3 logo">

[![publish status](https://github.com/blobtrtl3/trtl3/actions/workflows/publish.yml/badge.svg?branch=main)](https://github.com/blobtrtl3/trtl3/actions/workflows/publish.yml)
[![dockerhub](https://img.shields.io/docker/pulls/nothiaki/trtl3.svg)](https://hub.docker.com/r/nothiaki/trtl3)
[![go report](https://goreportcard.com/badge/github.com/blobtrtl3/trtl3)](https://goreportcard.com/report/github.com/gin-gonic/gin)
[![APGL-3.0](https://img.shields.io/badge/license-AGPL--3.0-blue.svg)](LICENSE)
[![PRs welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](CONTRIBUTING.md)

**Trtl3**🐢 is a file storage service written in [GO](https://go.dev/) designed for developers who want to build things without
the complexity of cloud services or external dependencies.
If you're building something with files uploads, downloads, and blobs organization — Trtl3 is a great starting point!

---

## Table of Contents

- [Getting Started](#getting-started)
- [Usage](#usage)
- [Features](#features)

---

## Getting Started

After Install Docker/Docker Compose you can install trtl3 using this command:

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
```

---

## Usage

So you can use it with SDK's or REST.
Go on our [website](https://trtl3.store) to learn more.

---

## Features

soon...
