---
date: 2026-07-10
priority: 2
status: filed
label: Support tickets need file attachments and a structured API for agent-driven bug reports
---

Filed through Audos Help directly, prompted by the support-chat formatting gap in `bugs/0036`: a
technical DNS/API evidence reply got flattened into an unreadable block, with no way to attach a
screenshot, log file, or HAR/JSON instead of pasting raw text into a chat box.

The filed request asks Audos to:

- add file attachments to Help/Priority Support tickets and replies — screenshots, markdown/text, logs,
  HAR files, JSON, CSV, small ZIPs;
- preserve filename, size, upload timestamp, uploader identity, and downloadability;
- expose a structured ticket API — create, read, reply, attach — so an agent (not just a human in the
  chat UI) can file a properly-structured report with workspace ID, type, title, markdown description,
  priority, source URL, environment metadata, repro steps, expected/actual behavior, and evidence
  attachments as distinct fields, not one paragraph;
- keep the same workspace auth/security model and avoid exposing secrets in ordinary chat logs by default.

Source: `throughline/docs/working/audos-otto-browser-approach.md`.
