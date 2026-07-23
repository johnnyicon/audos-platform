---
date: 2026-07-13
area: app-build
status: confirmed
label: Does pinning cursorModel actually change which model runs, or is it ignored?
---

**Hypothesis:** Explicitly pinning a `cursorModel` (e.g. `fable-5`) on a dispatched job changes which
model actually executes it, rather than being an accepted-but-ignored parameter.

**Method:** Dispatched a throwaway app briefed to self-report which model it believed it was running as,
then read the live rendered page directly rather than trusting the job's (often-truncated) completion
report.

**Result: confirmed pass.** The live page reported "Fable 5" / internal slug `claude-fable-5-thinking`
— appropriately hedged in our own notes as system-prompt-sourced, not something the model can fully
verify about itself from the inside. A later single-run reuse of the same pin on a real bug fix is
explicitly *not* counted as a second confirmation — one successful run doesn't prove cause-and-effect on
its own.

See `docs/platform/21-cursor-backend-research.md`, "Model pinning — confirmed working, verified live."
