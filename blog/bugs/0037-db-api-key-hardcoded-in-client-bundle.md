---
date: 2026-04-20
area: db-api
status: open
filed: no
label: A workspace's DB_API_KEY was found hardcoded in source that ships to the browser — fully exposed in the compiled bundle
---

`DB_API_KEY` — a secret used for direct workspace-DB-hook access — was found hardcoded in
`audos-config.ts`, client-side source that compiles straight into the browser bundle for any Audos
CDN-deployed app using direct DB-hook access from React. Anyone who opens devtools and reads the
deployed JavaScript has the key in plaintext. Not something introduced by porting the app elsewhere —
the exposure already existed in the Audos-hosted deployment; porting just carried it forward unchanged.

Flagged rather than fixed on the spot, since resolving it properly means moving direct DB-hook calls
that need this key behind a server-side proxy instead of calling from the browser at all — a real
architecture change, not a config tweak. Worth treating as a standing warning for anyone building a
client-side Audos app that touches `db-api` hooks directly: don't assume secrets referenced from React
component code stay server-side just because the docs describe them as "API keys."

Source: `throughline/docs/working/audos-migration-plan.md`.
