# FRONTEND TOOLING KNOWLEDGE BASE

## OVERVIEW
`web/` is the Vue/Vite toolchain boundary: scripts, env files, build config, Tailwind theme, lint/format rules, and the `src/` app subtree.

## STRUCTURE
```text
web/
├── package.json
├── scripts/setup.js
├── vite.config.ts
├── tailwind.config.js
├── .eslintrc.cjs
├── .prettierrc
├── .env*
└── src/            # actual app code (see child AGENTS)
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Dev/build/lint commands | `package.json` | authoritative frontend scripts |
| First-time setup | `scripts/setup.js` | intended Node/version/typecheck flow; currently broken by ESM/CJS mismatch |
| Dev server + aliases + proxy | `vite.config.ts` | port `3000`, `/api` proxy to backend, import aliases |
| Theme tokens + motion | `tailwind.config.js` | glassmorphism colors, dark mode, custom animations |
| Lint rules | `.eslintrc.cjs` | Vue + TS rules, explicit emits, warning levels |
| Formatting rules | `.prettierrc` | no semicolons, single quotes, print width 100 |
| App code | `src/` | child `AGENTS.md` is more specific |

## CONVENTIONS
- Node 18+ is the expected baseline.
- `scripts/setup.js` is intended to prefer Yarn when available, otherwise npm, but current verification shows it crashes under Node because this package is ESM and the script still uses `require`.
- Production build is `vue-tsc && vite build`, not just `vite build`.
- Tailwind is the design system source of truth (`darkMode: 'class'`, custom glass/background/text palettes, custom animations).
- Import aliases are configured in Vite/TypeScript; use them instead of long relative paths.

## ANTI-PATTERNS
- Do not put real application logic at the `web/` root; that belongs in `web/src/`.
- Do not invent a parallel styling system outside Tailwind tokens and existing assets/styles.
- Do not assume a frontend test runner exists; `package.json` does not define a `test` script.
- Do not assume `scripts/setup.js` is currently a working bootstrap path; use `npm ci`/`npm install` until the ESM/CommonJS mismatch is fixed.

## NOTES
- `npm run lint` uses `--fix`, so it may rewrite files.
- Current repo reality: `node scripts/setup.js` fails under Node 24 because `web/package.json` declares `"type": "module"` while the script still uses CommonJS `require(...)`.
- Use `web/src/AGENTS.md` for route/state/view/component conventions; this file is intentionally tooling-focused.
