# MAILSERVER KNOWLEDGE BASE

## OVERVIEW
SMTP receive, SMTP submit, and IMAP runtime live here; this subtree is protocol/server code, not HTTP API code.

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Server lifecycle + ports | `server.go` | builds shared storage, starts/stops all protocol servers |
| SMTP receive/submit behavior | `smtp_backend.go`, `smtp_session.go` | inbound vs authenticated submission flow |
| IMAP wiring + sessions | `imap.go`, `imap_session.go`, `imap_session_selected.go` | session states and selected-mailbox behavior |
| Storage bridge | `storage.go` | shared persistence access used by protocol servers |

## CONVENTIONS
- Local config expects separate ports for SMTP receive (`25`), SMTP submit (`587`), and IMAP TLS (`993`).
- `NewMailServer(...)` creates one shared storage layer and then wires all protocol servers from that shared base.
- `Start()`/`Stop()` own goroutine lifecycle. `main.go` is responsible for starting this subsystem alongside the HTTP app.
- Keep external inbound SMTP behavior separate from authenticated user submission behavior.
- TLS certificate/key paths come from config; they are not hardcoded in this package.

## ANTI-PATTERNS
- Do not move HTTP handler, router, or frontend concerns into this subtree.
- Do not hardcode ports, certificate paths, or domains here.
- Do not reimplement auth/JWT helpers in protocol code; shared auth utilities stay in `pkg/auth/` or middleware layers.
- Do not assume all security checks are finished; `smtp_session.go` still contains TODOs around sender authorization and spam/virus checks.

## NOTES
- `server.go` logs explicit operator guidance about the difference between external mail reception and user submission ports.
- Port testing code in startup is present but commented/optional, so startup success does not prove every listener is reachable externally.
- For config-shape changes, update `internal/config/config.go` first and keep this subtree aligned with it.
