<img alt="Trtl3 Logo" src="https://github.com/blobtrtl3/docs/blob/main/logo/logo.svg" width="200"/>

[![publish status](https://github.com/blobtrtl3/trtl3/actions/workflows/publish.yml/badge.svg?branch=main)](https://github.com/blobtrtl3/trtl3/actions/workflows/publish.yml)
[![dockerhub](https://img.shields.io/docker/pulls/nothiaki/trtl3.svg)](https://hub.docker.com/r/nothiaki/trtl3)
[![go report](https://goreportcard.com/badge/github.com/blobtrtl3/trtl3)](https://goreportcard.com/report/github.com/gin-gonic/gin)
[![APGL-3.0](https://img.shields.io/badge/license-AGPL--3.0-blue.svg)](LICENSE)
[![PRs welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](CONTRIBUTING.md)

## Trtl3

**Trtl3**🐢 is a file storage service designed for beginner developers who want to learn how object storage works in practice — without
the complexity of cloud services or external dependencies.

If you're building or experimenting with file uploads, downloads, and basic file organization via HTTP — Trtl3 is a great starting point!

---

## 🛠️ Getting Started

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

So you can use it with SDK's or REST.
Go on [Trtl3 Website](https://trtl3.store) to know more.

