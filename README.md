# lite-gitops

A lightweight GitOps agent written in Go.

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
