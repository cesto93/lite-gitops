# lite-gitops agent

This project is a lightweight GitOps agent in Go.

## Deploy

Docker Compose is supported natively. When `compose.files` is configured in the YAML config,
the `reconcile()` method runs `docker compose pull && docker compose up -d` against the
synced repo directory. The compose binary must be available on the system PATH.

## Build

```bash
go build -o lite-gitops .
```

## Related

- `pkg/compose/compose.go` — Docker Compose deployment logic
- `pkg/agent/agent.go` — Agent loop and reconcile entry point
- `pkg/config/config.go` — Config struct including ComposeConfig
