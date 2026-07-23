---
date: 2026-07-08
priority: 2
status: filed
label: DNS for an Audos-managed custom domain requires a support ticket — no self-serve write access, even via Otto's own API
---

Every DNS change on an Audos-managed custom domain — even a single CNAME repoint or TXT record removal
— has to go through a Human-In-The-Loop support ticket. Otto's own DNS tools are read-only; there is no
API or UI path for a developer (or an agent acting on their behalf) to make the change directly. In
practice this turned a same-day DNS cutover into a multi-day wait, including one case where a support
reply marked a ticket "done" while the live DNS record hadn't actually changed — confirmed by checking
public DNS resolution directly rather than trusting the ticket status.

The investigation also mapped the actual mechanism, worth documenting as a capability boundary rather
than just a complaint: Audos's custom-domain hosting is built on Cloudflare-for-SaaS custom-hostname
claims — a `_cf-custom-hostname.<subdomain>` TXT record marking the claim, paired with an
`_acme-challenge.<subdomain>` TXT for ACME validation. Releasing a subdomain to point elsewhere (e.g. to
Railway) means explicitly deleting both records during cutover — there's no separate "unbind" tool, DNS
record removal *is* the release mechanism.

The ask: self-serve write access to DNS for domains a workspace owns, without a ticket round-trip for
routine changes — with the Cloudflare-for-SaaS release sequence above exposed as a documented, supported
operation rather than something a support engineer has to perform by hand each time.

Source: `throughline/docs/working/audos-otto-browser-approach.md`.
