# lite-gitops agent

This project is a lightweight GitOps agent in Go.

## Deploy

Docker Compose is supported natively. When `compose.files` is configured in the YAML config,
the `reconcile()` method runs `docker compose pull && docker compose up -d` against the
synced repo directory. The compose binary must be available on the system PATH.

## Build

Module path is `github.com/cesto93/lite-gitops`, so internal imports use that prefix.

```bash
go build -o lite-gitops .
go install .   # installs to $(go env GOPATH)/bin
```

## Container

`Dockerfile` produces a minimal Alpine image via multi-stage Go build.
The container bundles docker-cli, docker-compose, git, and openssh so it
can clone repos and run Compose on the host Docker socket.

`docker-compose.yml` runs the agent with the Docker socket mounted,
a read-only config mount at `/etc/lite-gitops/config.yaml`, and a
persistent volume for the Git clone state.

`.github/workflows/docker-build.yml` builds the image on every push to
`main` or tag `v*` and saves the image tarball as a workflow artifact.

## Related

- `pkg/compose/compose.go` — Docker Compose deployment logic
- `pkg/agent/agent.go` — Agent loop and reconcile entry point
- `pkg/config/config.go` — Config struct including ComposeConfig
