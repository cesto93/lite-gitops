# lite-gitops

A lightweight GitOps agent written in Go.

## Install

Install straight from the repository (Go 1.26.5 or newer, as declared in `go.mod`):

```bash
go install github.com/cesto93/lite-gitops@latest
```

Or install from a local clone:

```bash
git clone https://github.com/cesto93/lite-gitops
cd lite-gitops
go install .
```

Either way the binary lands in `$(go env GOPATH)/bin`. Make sure that directory is
on your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Docker

Prebuilt multi-arch images (`linux/amd64`, `linux/arm64`) are published to the GitHub
Container Registry on every push to `main` and on `v*` tags:

```bash
docker pull ghcr.io/cesto93/lite-gitops:main
```

The container needs the host Docker socket to deploy Compose projects and a config file
mounted read-only:

```bash
docker run -d --name lite-gitops --restart unless-stopped \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ./config.yaml:/etc/lite-gitops/config.yaml:ro \
  -v lite-gitops-data:/var/lib/lite-gitops \
  ghcr.io/cesto93/lite-gitops:main run --config /etc/lite-gitops/config.yaml
```

Or with Compose, using the `docker-compose.yml` in this repository:

```yaml
services:
  lite-gitops:
    image: ghcr.io/cesto93/lite-gitops:main
```

## Usage

```
lite-gitops run --config /etc/lite-gitops/config.yaml
```

### Flags

| Flag          | Default                              | Description              |
|---------------|--------------------------------------|--------------------------|
| `--config`    | `~/.lite-gitops/config.yaml`         | Path to config file      |
| `-v`/`--verbose` | `false`                           | Enable debug logging     |

## Config

```yaml
git:
  repository: https://github.com/org/my-configs.git  # required
  branch: main                                        # default: main
  remote: origin                                      # default: origin
  local_path: /var/lib/lite-gitops/state              # clone destination
compose:
  files:                                              # default: [docker-compose.yml]
    - docker-compose.yml
  project_name: my-project                            # optional
  build: false                                        # default: false
poll_interval: 30s                                    # default: 30s
verbose: false
```

## Docker Compose

When configured, `lite-gitops` automatically deploys Docker Compose projects from the
synced Git repository. On every new commit, it runs `docker compose pull` followed by
`docker compose up -d` to apply changes.

## Build

```bash
go build -o lite-gitops .
```
