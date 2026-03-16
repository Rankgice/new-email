# BACKEND KNOWLEDGE BASE

## OVERVIEW
Layered Go backend for HTTP APIs, persistence, service orchestration, auth, and mail runtime integration.

## STRUCTURE
```text
internal/
├── config/      # YAML config schema + loading
├── handler/     # HTTP handlers
├── middleware/  # auth/admin/API-key guards
├── model/       # DB-facing models
├── result/      # response helpers
├── router/      # route registration
├── service/     # SMTP/IMAP/SMS/storage/cache orchestration
├── svc/         # dependency container
├── types/       # request/response DTOs
└── mailserver/  # protocol subsystem (see child AGENTS)
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Add/change endpoint | `router/router.go` + `handler/` + `types/` | routes, handler wiring, payload shapes |
| Auth/middleware changes | `middleware/` + `pkg/auth/` | middleware stays here, token/password helpers live in `pkg/auth/` |
| DB-backed behavior | `model/` | keep SQL concerns out of handlers |
| Shared dependencies | `svc/service_context.go` | DB, config, models, service manager, MinIO |
| External service orchestration | `service/` | SMTP/IMAP/SMS/storage/cache manager + connection tests |
| DTO/schema edits | `types/` | request/response structs + validation tags |
| Mail protocol runtime | `mailserver/` | separate child rules apply there |

## CONVENTIONS
- `ServiceContext` is the backend center of gravity. New handlers/services should consume dependencies from it instead of rebuilding clients/config.
- Typical handler flow here is: bind request -> derive current user/admin/API key from middleware -> validate ownership/existence -> call model/service -> respond through `internal/result`.
- `internal/types/` is the home for API payload structs and validation-tagged request/response shapes.
- `internal/service/` is for integration orchestration and connectivity checks, not for route registration or direct HTTP response shaping.
- Route changes are rarely one-file edits: expect to touch `router`, `handler`, `types`, and often `model` or `service` together.

## ANTI-PATTERNS
- No raw SQL or hand-built queries in handlers/services; architecture docs explicitly ban this.
- Do not put DB work directly in handlers.
- Do not bypass `internal/result` with ad hoc JSON shapes unless a local pattern already does so for that endpoint family.
- Do not call config loaders directly from random backend packages; use config carried through `ServiceContext`.
- Do not move protocol-server logic into generic handlers/services; mail protocol behavior belongs in `internal/mailserver/`.

## NOTES
- `internal/handler/` is central and partly incomplete: mailbox sync, domain verification, and API-user resolution still have TODO-heavy areas.
- `internal/types/` is dense but uniform; prefer extending existing DTO patterns instead of inventing one-off structs in handlers.
- Shared JWT/password semantics live in `pkg/auth/`, even though that package is outside `internal/`.
