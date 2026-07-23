---
date: 2026-07-14
area: storage-api
status: open
filed: no
label: "`/api/upload/image` silently corrupts non-image payloads"
---

A base64-encoded PDF uploaded through the image-upload endpoint returns `HTTP 200 {"success":true}` with
a plausible-looking `.png` URL — but the stored bytes don't match what was sent.

> Confirmed via md5: sent hash `7fc99347...`, stored hash `ec78847b...`, 3,031 bytes sent vs. 3,050 bytes
> stored, zero common prefix — genuinely corrupted, not just re-encoded. A real `image/png` payload
> round-trips md5-identical through the same endpoint; the multipart endpoint (`/api/upload/file`)
> round-trips correctly for any content type. The danger here isn't the corruption itself — it's that it
> fails **silently**, with a success response and a working-looking URL. Any non-image file (PDF, audio,
> video) must go through `/api/upload/file`, never `/api/upload/image`, regardless of what the endpoint's
> name might suggest is fine.
