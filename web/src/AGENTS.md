# FRONTEND APP KNOWLEDGE BASE

## OVERVIEW
`web/src/` contains the actual Vue application: bootstrap, router, state, views, components, API wrappers, shared types, and UI utilities.

## STRUCTURE
```text
web/src/
├── main.ts
├── router/
├── stores/
├── views/
├── components/
├── api/
├── composables/
├── types/
└── utils/
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| App bootstrap/plugins | `main.ts` | auth init, router, Vue Query, motion, global error handler |
| Navigation/guards/layouts | `router/index.ts` | route meta, lazy imports, auth/admin gating |
| Global state | `stores/auth.ts`, `stores/theme.ts` | only two Pinia stores right now |
| Route-level screens | `views/` | grouped by `admin`, `auth`, `email`, `error`, `user`, `verification-codes` |
| Shared UI | `components/` | `ui`, `email`, `mailbox`, `settings` groups |
| Frontend API wrappers | `api/` | thin modules for admin/settings/verification-code flows |
| Reusable hooks | `composables/` | small shared helpers like theme/notification/responsive |

## CONVENTIONS
- `main.ts` initializes auth before mount and registers router, Vue Query, and VueUse Motion.
- Route metadata is meaningful here: `requiresAuth`, `requiresAdmin`, `layout`, `title`, and `keepAlive` affect behavior.
- Routes are lazy-loaded and grouped by product area; follow the existing auth/email/user/admin split.
- Stores stay small and central; shared logic belongs in composables only when it is truly cross-view.
- Route-level screens belong in `views/`; supporting building blocks belong in `components/`.
- Use configured aliases rather than deep relative imports.

## ANTI-PATTERNS
- Do not duplicate auth/role gating in random components when router meta or the auth store should own it.
- Do not move theme tokens into feature files; `web/tailwind.config.js` and shared styles own the visual system.
- Do not create deeper AGENTS files yet; this subtree is still coherent enough to document as one unit.
- Do not add new pages without updating route meta intentionally.

## NOTES
- `views/verification-codes/components/` is the deepest nested frontend area; check for existing local patterns before creating new ones.
- Admin/settings surfaces already have specialized components—reuse them before introducing one-off patterns.
