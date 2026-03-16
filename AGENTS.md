# PROJECT KNOWLEDGE BASE

**Generated:** 2026-03-08
**Repo state:** not a git repo

## OVERVIEW
Go email platform with a Gin/GORM backend, YAML runtime config, built-in SMTP/IMAP services, and a Vue 3 + Vite frontend for user/admin workflows.

## STRUCTURE
```text
./
├── main.go                  # backend + mailserver startup
├── etc/config.yaml          # active runtime config
├── internal/                # backend application code
│   └── mailserver/          # protocol-specific SMTP/IMAP runtime
├── pkg/auth/                # JWT + password helpers
├── web/                     # frontend tooling/config boundary
│   └── src/                 # actual SPA code
└── *.md docs                # API, deployment, architecture, planning notes
```

## HIERARCHY
- `./AGENTS.md` = repo-wide defaults
- `internal/AGENTS.md` = backend layering and API rules
- `internal/mailserver/AGENTS.md` = protocol/server rules
- `web/AGENTS.md` = frontend tooling/config rules
- `web/src/AGENTS.md` = frontend application-code rules

Nearest `AGENTS.md` wins.

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Runtime config schema | `internal/config/config.go` | YAML-backed config model for app/web/db/jwt/mail/redis/minio |
| Backend startup | `main.go` | builds config, `ServiceContext`, router, and mailserver |
| HTTP route map | `internal/router/router.go` | public/user/admin/API-key route groups |
| Backend dependencies | `internal/svc/service_context.go` | DB, models, MinIO, service manager |
| HTTP handlers | `internal/handler/` | request binding + auth context + result responses |
| External integrations | `internal/service/` | SMTP/IMAP/SMS/storage/cache orchestration |
| Protocol servers | `internal/mailserver/` | SMTP receive/submit + IMAP runtime |
| Auth helpers | `pkg/auth/` | JWT generation/parsing + password hashing |
| Frontend tooling | `web/package.json`, `web/vite.config.ts`, `web/tailwind.config.js` | scripts, aliases, proxy, theme tokens |
| Frontend bootstrap | `web/src/main.ts` | auth init + router + Vue Query + motion |
| Frontend navigation | `web/src/router/index.ts` | lazy routes, meta, layouts, auth/admin gating |

## CONVENTIONS
- Runtime config is centralized in `etc/config.yaml`; code changes should follow the schema in `internal/config/config.go`.
- Backend startup path is `main.go` -> `svc.NewServiceContext(...)` -> `router.SetupRouter(...)` plus `mailserver.NewMailServer(...)`.
- Frontend tooling lives in `web/`; frontend app code starts in `web/src/`.
- Vite dev server runs on port `3000` and proxies `/api` to `http://localhost:8080`.
- UI styling is anchored in `web/tailwind.config.js` (`darkMode: 'class'`, custom glass/background/text palettes, custom animations).
- Frontend expects Node 18+; `web/scripts/setup.js` is intended to enforce setup/typecheck, but it currently fails under Node because the package is ESM while the script still uses `require`.

## ANTI-PATTERNS (THIS PROJECT)
- No raw/direct SQL. Architecture docs explicitly forbid it.
- Do not put DB access in handlers; use models/services via `ServiceContext`.
- Do not reload config ad hoc from random packages; consume config passed through context.
- Do not expose API keys in UI, docs, or examples.
- Do not assume README feature claims mean every path is production-complete; several handlers/mailserver paths still contain TODO placeholders.
- Do not claim `etc/config.yaml.example` exists unless you add it; only `etc/config.yaml` is present right now.

## UNIQUE STYLES
- Backend responses are normalized through `internal/result`.
- Frontend route meta (`requiresAuth`, `requiresAdmin`, `layout`, `title`, `keepAlive`) is a real control surface, not decorative metadata.
- Auth helpers support JWT plus mixed password verification (`bcrypt` and `Argon2id` compatibility).

## COMMANDS
```bash
go run main.go -f etc/config.yaml
go build -o email-system main.go
go test ./...
cd web && npm ci
cd web && npm run dev
cd web && npm run build
cd web && npm run lint
docker compose up -d
```

## NOTES
- `docker-compose.yml` runs the app plus MinIO and exposes `8080`, `25`, `587`, and `993`.
- GitHub Actions currently focus on Docker image publishing, not a full lint/test matrix.
- The README still references `go test ./test/... -v` and `node scripts/setup.js`, but current verification showed the test directory does not exist and the setup script crashes under the repo's ESM config.
- If you are working inside `internal/`, `internal/mailserver/`, `web/`, or `web/src/`, prefer the closer `AGENTS.md` over this file.
