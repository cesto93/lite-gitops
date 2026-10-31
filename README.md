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
  username: my-github-user                            # optional, private repos
  token: ghp_xxxxxxxxxxxxxxxxxxxx                     # optional, private repos
compose:
  files:                                              # default: [docker-compose.yml]
    - docker-compose.yml
  project_name: my-project                            # optional
  build: false                                        # default: false
poll_interval: 30s                                    # default: 30s
verbose: false
```

### Private repositories

Cloning and pulling private GitHub repositories over HTTPS works with a personal
access token (classic or fine-grained). Set `git.username` to your GitHub username
and `git.token` to the token. If `git.username` is omitted, `x-access-token` is used,
which is the convention for GitHub App installation tokens.

The repository URL must use HTTPS, not SSH:

```yaml
git:
  repository: https://github.com/org/private-configs.git
  username: my-github-user
  token: ghp_xxxxxxxxxxxxxxxxxxxx
```

Prefer providing the token through the environment instead of committing it to the
config file. Every option can be overridden with an environment variable prefixed by
`LITE_GITOPS_`, using underscores instead of dots:

```bash
export LITE_GITOPS_GIT_USERNAME=my-github-user
export LITE_GITOPS_GIT_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx
lite-gitops run --config /etc/lite-gitops/config.yaml
```

With Docker:

```bash
docker run -d --name lite-gitops --restart unless-stopped \
  -e LITE_GITOPS_GIT_USERNAME=my-github-user \
  -e LITE_GITOPS_GIT_TOKEN=ghp_xxxxxxxxxxxxxxxxxxxx \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v ./config.yaml:/etc/lite-gitops/config.yaml:ro \
  -v lite-gitops-data:/var/lib/lite-gitops \
  ghcr.io/cesto93/lite-gitops:main run --config /etc/lite-gitops/config.yaml
```

Or with Compose, export the two variables in your shell before running `docker compose up`;
`docker-compose.yml` forwards them to the container when present.

The token only needs read access to the repository contents. The agent never logs it.

## Docker Compose

When configured, `lite-gitops` automatically deploys Docker Compose projects from the
synced Git repository. On every new commit, it runs `docker compose pull` followed by
`docker compose up -d` to apply changes.

## Build

```bash
go build -o lite-gitops .
```
