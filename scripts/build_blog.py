#!/usr/bin/env python3
"""
Build the Audos SDK build-log page from the Markdown posts in blog/.

Usage: python3 scripts/build_blog.py
Output: blog/dist/index.html — a single self-contained HTML file (font embedded),
        safe to open directly in a browser or publish as a Claude Artifact.

The page has four tab sections sharing one masthead/nav (all in one file, since
each Claude Artifact publish is a single URL): Blog, Experiments, Bugs, Feature Requests.

Adding a new blog post: drop a new blog/NNNN-slug.md file with the frontmatter block
shown below, update blog/INDEX.md, and re-run this script. See blog/HOW-TO-UPDATE.md.

Adding a new experiment: drop a new blog/experiments/NNNN-slug.md file (frontmatter:
date, area, status: confirmed|corrected|inconclusive|open, label). An experiment is a
hypothesis-driven test of a specific capability — distinct from a bug (broken) or a
feature request (missing).

Adding a new bug: drop a new blog/bugs/NNNN-slug.md file (frontmatter: date, area,
status: fixed|open, filed: yes|no, filed_ref: <path, if filed>, label).

Adding a new feature request: drop a new blog/feature-requests/NNNN-slug.md file
(frontmatter: date, priority: 1|2|3, status: "not filed"|filed, label).

No third-party dependencies — stdlib only, so any agent (Claude, Codex, etc.) can run
this with a bare `python3`.
"""
import base64
import html
import re
from pathlib import Path

REPO_ROOT = Path(__file__).resolve().parent.parent
BLOG_DIR = REPO_ROOT / "blog"
BUGS_DIR = BLOG_DIR / "bugs"
FR_DIR = BLOG_DIR / "feature-requests"
EXPERIMENTS_DIR = BLOG_DIR / "experiments"
PLATFORM_DIR = REPO_ROOT / "docs" / "platform"
FONT_PATH = BLOG_DIR / "assets" / "jetbrains-mono-var.woff2"
OUT_PATH = BLOG_DIR / "dist" / "index.html"

# One-line blurbs for the SDK tab — hand-curated (titles are read live from each
# file, so they can't drift; blurbs need a one-line update here when a new numbered
# doc lands in docs/platform/, same as skill/SKILL.md's own lookup table).
PLATFORM_DOC_BLURBS = {
    "01-qa-answers.md": "Answers to specific platform questions raised during development, gathered in one place.",
    "02-platform-overview.md": "What Audos is and how the top-level pieces fit together.",
    "03-workspace-structure.md": "How a workspace's folders and files are actually organized.",
    "04-development-workflow.md": "When to develop locally vs. on-platform, and how to make both work together.",
    "05-composite-apis.md": "How to combine multiple Audos primitives into one higher-level server function.",
    "06-capabilities-reference.md": "A full capability reference — corrected in place after a live probe found several original claims wrong.",
    "07-api-development-guide.md": "How to build external API endpoints (server functions) on Audos.",
    "08-server-function-runtime.md": "The sandboxed JS runtime server functions actually run in — not full Node.js.",
    "09-server-function-templates.md": "Copy-paste server function templates for common patterns (CRUD, AI, email...).",
    "10-platform-architecture.md": "The stack underneath Audos apps: React, Tailwind, hooks, and how they fit together.",
    "11-typescript-api-client.md": "The REST API surface — base URL pattern and how to call hooks from anywhere.",
    "12-local-development.md": "A local mock layer that mirrors the real platform so local dev behaves identically.",
    "13-database-access.md": "Authoritative reference for the workspace Postgres schema — creation, ownership, gotchas.",
    "14-workspace-db-analysis.md": "A point-in-time full analysis of what's actually in a live workspace database.",
    "15-workspace-db-table-ownership.md": "Classifies every workspace table — what to keep, drop, or ask about.",
    "16-github-dev-mode-app-deployment.md": "How the GitHub Dev Mode deploy pipeline works, and how to debug a blank page fast.",
    "17-emailgate-otp-configuration.md": "How to configure OTP verification on EmailGate — read doc 26 first if you have custom auth.",
    "18-genesis-space-and-ui-ceiling.md": "The line between “workspace” and “app,” and how far custom UI can actually go.",
    "19-capabilities-and-limitations.md": "Verified, dated platform limits — re-check before relying on any of them, Audos moves fast.",
    "20-power-user-wishlist.md": "What a real developer/power-user channel on Audos would need, grounded in specific incidents.",
    "21-cursor-backend-research.md": "Confirms Audos's App Agent is genuinely Cursor's Background Agent product underneath.",
    "22-eliminating-the-chat-shell-playbook.md": "Step-by-step recipe for a real full-screen app experience instead of the default chat shell.",
    "23-dispatch-then-poll-sop.md": "The standard operating pattern for every job dispatched to Otto/Cursor — dispatch, then poll.",
    "24-where-new-findings-go.md": "The decision rule for where a new finding actually belongs: log, backlog, doc, or blog.",
    "25-escalation-and-support-paths.md": "The two real escalation mechanisms that exist for a genuinely stuck job, verified live.",
    "26-unified-space-signed-out-view-is-mandatory.md": "Why EmailGate can't just be disabled, and what happened when we tried to replace it.",
    "27-blog-hosting-cloudflare-worker.md": "How this very blog is hosted — Cloudflare Worker, Basic Auth, and the mistakes made getting there.",
}

CHIP_CLASS = {
    "fixed": "chip-fixed",
    "shipped": "chip-fixed",
    "pass": "chip-pass",
    "open": "chip-open",
    "filed": "chip-fixed",
    "not filed": "chip-open",
    "confirmed": "chip-pass",
    "corrected": "chip-fixed",
    "inconclusive": "chip-open",
}


def parse_frontmatter(path: Path) -> tuple[dict, str]:
    text = path.read_text(encoding="utf-8")
    m = re.match(r"^---\n(.*?)\n---\n(.*)$", text, re.DOTALL)
    if not m:
        raise ValueError(f"{path.name}: missing --- frontmatter block")
    fm_raw, body = m.groups()
    meta = {}
    for line in fm_raw.splitlines():
        if ":" in line:
            k, v = line.split(":", 1)
            meta[k.strip()] = v.strip().strip('"')
    return meta, body.strip()


def inline_md(text: str) -> str:
    text = html.escape(text, quote=False)
    text = re.sub(r"\*\*(.+?)\*\*", r"<strong>\1</strong>", text)
    text = re.sub(r"(?<!\*)\*(?!\*)(.+?)(?<!\*)\*(?!\*)", r"<em>\1</em>", text)
    text = re.sub(r"`([^`]+)`", r"<code>\1</code>", text)
    return text


def render_paragraph(p: str) -> str:
    is_callout = p.startswith(">")
    if is_callout:
        p = " ".join(line.lstrip(">").strip() for line in p.splitlines())
        cls = "callout"
    else:
        p = " ".join(line.strip() for line in p.splitlines())
        cls = None
    inner = inline_md(p)
    return f'<p class="{cls}">{inner}</p>' if cls else f"<p>{inner}</p>"


def render_body(body: str) -> str:
    paragraphs = [p.strip() for p in re.split(r"\n\s*\n", body) if p.strip()]
    return "\n        ".join(render_paragraph(p) for p in paragraphs)


# ─── Blog posts ─────────────────────────────────────────────────────────────

def parse_post(path: Path) -> dict:
    meta, body = parse_frontmatter(path)
    title_m = re.search(r"^#\s+(.+)$", body, re.MULTILINE)
    meta["title"] = title_m.group(1).strip() if title_m else path.stem
    body = re.sub(r"^#\s+.+\n", "", body, count=1, flags=re.MULTILINE).strip()
    meta["body_html"] = render_body(body)
    return meta


def build_blog_section(posts: list[Path]) -> dict:
    entries = [parse_post(p) for p in posts]
    scorecard_cells, index_links, articles = [], [], []
    for i, e in enumerate(entries, start=1):
        eid = f"entry-{i}"
        chip = CHIP_CLASS.get(e.get("status", ""), "chip-open")
        num = f"{i:02d}"
        scorecard_cells.append(f"""
    <div class="score-cell" onclick="document.getElementById('{eid}').scrollIntoView({{behavior:'smooth'}})">
      <div class="score-num"><span>{num}</span><span class="prod">{html.escape(e['product'])}</span></div>
      <div class="score-label">{html.escape(e['label'])}</div>
      <span class="chip {chip}">{html.escape(e['status'])}</span>
    </div>""")
        index_links.append(
            f'      <a href="#{eid}" data-target="{eid}">{num} &middot; {html.escape(e["label"])}'
            f'<span class="p">{html.escape(e["product"])}</span></a>'
        )
        articles.append(f"""
      <article id="{eid}">
        <div class="eyebrow"><span class="num">Entry {num}</span><span>{html.escape(e['date'])}</span><span>{html.escape(e['product'])}</span><span class="chip {chip}">{html.escape(e['status'])}</span></div>
        <h2>{html.escape(e['title'])}</h2>
        {e['body_html']}
      </article>""")
    return {
        "count": len(entries),
        "date_range": f"{entries[0]['date']} &rarr; {entries[-1]['date']}",
        "scorecard": "".join(scorecard_cells),
        "index_nav": "\n".join(index_links),
        "articles": "".join(articles),
    }


# ─── Bugs ───────────────────────────────────────────────────────────────────

def parse_bug(path: Path) -> dict:
    meta, body = parse_frontmatter(path)
    meta["body_html"] = render_body(body)
    return meta


AUDOS_STATUS_LABEL = {
    "pending": "filed with Audos · awaiting response",
    "acknowledged": "filed with Audos · acknowledged",
    "fixed": "filed with Audos · fixed by Audos",
    "wontfix": "filed with Audos · won't fix",
}


def build_bugs_section(paths: list[Path]) -> dict:
    items = [parse_bug(p) for p in paths]
    open_count = sum(1 for i in items if i.get("status") == "open")
    fixed_count = len(items) - open_count
    cards = []
    for i, item in enumerate(items, start=1):
        status = item.get("status", "open")
        chip = CHIP_CLASS.get(status, "chip-open")
        filed = item.get("filed", "no") == "yes"
        audos_status = item.get("audos_status", "pending") if filed else None
        filed_badge = (
            f'<span class="filed-badge filed-yes">{html.escape(AUDOS_STATUS_LABEL.get(audos_status, "filed with Audos"))}</span>'
            if filed else
            '<span class="filed-badge filed-no">not filed</span>'
        )
        ref = item.get("filed_ref", "")
        ref_html = f'<div class="filed-ref"><code>{html.escape(ref)}</code></div>' if ref else ""
        cards.append(f"""
      <div class="item-card" id="bug-{i:04d}">
        <div class="item-eyebrow">
          <span class="num">Bug {i:04d}</span><span>{html.escape(item['date'])}</span>
          <span class="area-tag">{html.escape(item.get('area',''))}</span>
          <span class="chip {chip}">{html.escape(status)}</span>
          {filed_badge}
        </div>
        <h3>{html.escape(item['label'])}</h3>
        {ref_html}
        {item['body_html']}
      </div>""")
    return {
        "count": len(items),
        "open_count": open_count,
        "fixed_count": fixed_count,
        "cards": "".join(cards),
    }


# ─── Feature requests ───────────────────────────────────────────────────────

def parse_fr(path: Path) -> dict:
    meta, body = parse_frontmatter(path)
    meta["body_html"] = render_body(body)
    return meta


def build_fr_section(paths: list[Path]) -> dict:
    items = [parse_fr(p) for p in paths]
    cards = []
    for i, item in enumerate(items, start=1):
        status = item.get("status", "not filed")
        chip = CHIP_CLASS.get(status, "chip-open")
        priority = item.get("priority", "3")
        cards.append(f"""
      <div class="item-card" id="fr-{i:04d}">
        <div class="item-eyebrow">
          <span class="num">FR {i:04d}</span><span>{html.escape(item['date'])}</span>
          <span class="area-tag">priority {html.escape(priority)}</span>
          <span class="chip {chip}">{html.escape(status)}</span>
        </div>
        <h3>{html.escape(item['label'])}</h3>
        {item['body_html']}
      </div>""")
    return {"count": len(items), "cards": "".join(cards)}


# ─── Experiments ────────────────────────────────────────────────────────────
# A structured, hypothesis-driven test of a specific platform capability —
# distinct from a bug (something broken) or a feature request (something
# missing). status: confirmed | corrected | inconclusive | open.

def parse_experiment(path: Path) -> dict:
    meta, body = parse_frontmatter(path)
    meta["body_html"] = render_body(body)
    return meta


EXPERIMENT_STATUS_LABEL = {
    "confirmed": "confirmed",
    "corrected": "corrected after re-test",
    "inconclusive": "inconclusive",
    "open": "still open",
}


def build_experiments_section(paths: list[Path]) -> dict:
    items = [parse_experiment(p) for p in paths]
    cards = []
    for i, item in enumerate(items, start=1):
        status = item.get("status", "open")
        chip = CHIP_CLASS.get(status, "chip-open")
        status_label = EXPERIMENT_STATUS_LABEL.get(status, status)
        cards.append(f"""
      <div class="item-card" id="exp-{i:04d}">
        <div class="item-eyebrow">
          <span class="num">Exp {i:04d}</span><span>{html.escape(item['date'])}</span>
          <span class="area-tag">{html.escape(item.get('area',''))}</span>
          <span class="chip {chip}">{html.escape(status_label)}</span>
        </div>
        <h3>{html.escape(item['label'])}</h3>
        {item['body_html']}
      </div>""")
    return {"count": len(items), "cards": "".join(cards)}


# ─── SDK reference docs (docs/platform/) ───────────────────────────────────
# Not a log — an index of the durable, generic platform knowledge this project
# has accumulated. Titles are read live from each file; blurbs are curated above.

# Docs that exist in docs/platform/ but shouldn't show up in the blog's SDK index —
# either not Audos-platform knowledge (27, infra for this blog itself) or superseded/
# redundant with later docs (01, an early FAQ mostly folded into 02/05/13 since). The
# files themselves stay in the repo either way; this only controls the blog listing.
SDK_EXCLUDE = {
    "01-qa-answers.md",
    "27-blog-hosting-cloudflare-worker.md",
}


def build_sdk_section() -> dict:
    paths = sorted(PLATFORM_DIR.glob("[0-9][0-9]-*.md")) if PLATFORM_DIR.exists() else []
    paths = [p for p in paths if p.name not in SDK_EXCLUDE]
    cards = []
    for p in paths:
        text = p.read_text(encoding="utf-8")
        title_m = re.search(r"^#\s+(.+)$", text, re.MULTILINE)
        title = title_m.group(1).strip() if title_m else p.stem
        num = p.name.split("-", 1)[0]
        blurb = PLATFORM_DOC_BLURBS.get(p.name, "")
        cards.append(f"""
      <div class="item-card" id="sdk-{num}">
        <div class="item-eyebrow">
          <span class="num">Doc {num}</span>
        </div>
        <h3>{html.escape(title)}</h3>
        <p>{html.escape(blurb)}</p>
        <div class="filed-ref"><code>docs/platform/{html.escape(p.name)}</code></div>
      </div>""")
    return {"count": len(paths), "cards": "".join(cards)}


def main():
    posts = sorted(BLOG_DIR.glob("[0-9][0-9][0-9][0-9]-*.md"))
    if not posts:
        raise SystemExit("No blog/NNNN-*.md posts found")
    bugs = sorted(BUGS_DIR.glob("[0-9][0-9][0-9][0-9]-*.md"))
    frs = sorted(FR_DIR.glob("[0-9][0-9][0-9][0-9]-*.md"))
    experiments = sorted(EXPERIMENTS_DIR.glob("[0-9][0-9][0-9][0-9]-*.md")) if EXPERIMENTS_DIR.exists() else []

    blog = build_blog_section(posts)
    bugs_section = build_bugs_section(bugs) if bugs else None
    fr_section = build_fr_section(frs) if frs else None
    exp_section = build_experiments_section(experiments) if experiments else None
    sdk_section = build_sdk_section()

    font_b64 = base64.b64encode(FONT_PATH.read_bytes()).decode("ascii")

    page = TEMPLATE.format(
        font_b64=font_b64,
        blog_count=blog["count"],
        date_range=blog["date_range"],
        blog_scorecard=blog["scorecard"],
        blog_index_nav=blog["index_nav"],
        blog_articles=blog["articles"],
        bugs_count=bugs_section["count"] if bugs_section else 0,
        bugs_open=bugs_section["open_count"] if bugs_section else 0,
        bugs_fixed=bugs_section["fixed_count"] if bugs_section else 0,
        bugs_cards=bugs_section["cards"] if bugs_section else "",
        fr_count=fr_section["count"] if fr_section else 0,
        fr_cards=fr_section["cards"] if fr_section else "",
        exp_count=exp_section["count"] if exp_section else 0,
        exp_cards=exp_section["cards"] if exp_section else "",
        sdk_count=sdk_section["count"],
        sdk_cards=sdk_section["cards"],
    )

    OUT_PATH.parent.mkdir(parents=True, exist_ok=True)
    OUT_PATH.write_text(page, encoding="utf-8")
    total = (blog["count"] + (bugs_section["count"] if bugs_section else 0)
             + (fr_section["count"] if fr_section else 0) + (exp_section["count"] if exp_section else 0))
    print(f"Wrote {OUT_PATH} ({len(page):,} bytes, {blog['count']} blog entries, "
          f"{bugs_section['count'] if bugs_section else 0} bugs, "
          f"{fr_section['count'] if fr_section else 0} feature requests, "
          f"{exp_section['count'] if exp_section else 0} experiments, "
          f"{sdk_section['count']} SDK docs, {total} items total)")


TEMPLATE = r'''<style>
@font-face{{
  font-family:'JBM'; font-style:normal; font-weight:100 900; font-display:swap;
  src:url(data:font/woff2;base64,{font_b64}) format('woff2');
}}
:root{{
  --paper:#F1F4F2; --paper-raised:#FFFFFF; --ink:#14201B; --ink-soft:#4B5C55; --ink-faint:#7C8C85;
  --line:#D8E0DC; --line-strong:#BFCBC5; --accent:#0F6E56; --accent-soft:#DCEEE7; --accent-ink:#083F31;
  --open:#B5502E; --open-soft:#F5E3DA; --open-ink:#5C2A16; --code-bg:#E7ECE9;
}}
@media (prefers-color-scheme: dark){{
  :root{{
    --paper:#101815; --paper-raised:#16211D; --ink:#E7EDEA; --ink-soft:#9FB0A8; --ink-faint:#6C7C75;
    --line:#29352F; --line-strong:#3A4A43; --accent:#4FD9B4; --accent-soft:#17322A; --accent-ink:#BFF3E1;
    --open:#E28A5C; --open-soft:#2E1F16; --open-ink:#F5CBAF; --code-bg:#1C2721;
  }}
}}
:root[data-theme="dark"]{{
  --paper:#101815; --paper-raised:#16211D; --ink:#E7EDEA; --ink-soft:#9FB0A8; --ink-faint:#6C7C75;
  --line:#29352F; --line-strong:#3A4A43; --accent:#4FD9B4; --accent-soft:#17322A; --accent-ink:#BFF3E1;
  --open:#E28A5C; --open-soft:#2E1F16; --open-ink:#F5CBAF; --code-bg:#1C2721;
}}
:root[data-theme="light"]{{
  --paper:#F1F4F2; --paper-raised:#FFFFFF; --ink:#14201B; --ink-soft:#4B5C55; --ink-faint:#7C8C85;
  --line:#D8E0DC; --line-strong:#BFCBC5; --accent:#0F6E56; --accent-soft:#DCEEE7; --accent-ink:#083F31;
  --open:#B5502E; --open-soft:#F5E3DA; --open-ink:#5C2A16; --code-bg:#E7ECE9;
}}
*{{box-sizing:border-box}}
.bloglog{{background:var(--paper);color:var(--ink);font-family:Iowan Old Style, Palatino Linotype, Palatino, Georgia, serif;line-height:1.7;padding:2.5rem 1.25rem 4rem}}
.bloglog a{{color:var(--accent)}}
.bloglog h1,.bloglog h2,.bloglog h3{{font-family:'JBM', ui-monospace, SFMono-Regular, Menlo, monospace;font-weight:700;line-height:1.22;text-wrap:balance;color:var(--ink);margin:0}}
.wrap{{max-width:1160px;margin:0 auto}}
.masthead{{display:flex;justify-content:space-between;align-items:flex-end;gap:1.5rem;flex-wrap:wrap;border-bottom:1px solid var(--line-strong);padding-bottom:1.5rem;margin-bottom:1rem}}
.masthead-title{{font-size:clamp(1.4rem,3vw,1.9rem);letter-spacing:-0.01em}}
.masthead-title .dot{{color:var(--accent)}}
.masthead-sub{{font-size:0.95rem;color:var(--ink-soft);max-width:50ch;margin-top:0.5rem}}
.masthead-meta{{font-family:'JBM', ui-monospace, monospace;font-size:0.7rem;letter-spacing:0.06em;text-transform:uppercase;color:var(--ink-faint);text-align:right;line-height:1.9;font-variant-numeric:tabular-nums}}
.tabnav{{position:sticky;top:0;z-index:30;display:flex;align-items:center;gap:0.4rem;margin-bottom:2rem;background:var(--paper);border-bottom:1px solid var(--line);padding-bottom:0;transition:box-shadow 0.2s ease,border-color 0.2s ease}}
.tabnav.stuck{{box-shadow:0 2px 10px rgba(0,0,0,0.08);border-bottom-color:var(--line-strong)}}
.sticky-title{{display:inline-flex;align-items:center;font-family:'JBM', monospace;font-weight:700;font-size:0.82rem;line-height:1;color:var(--ink);white-space:nowrap;margin-right:0.9rem;padding:0.7rem 0;cursor:pointer;opacity:0;transform:translateY(-6px) scale(0.96);pointer-events:none;transition:opacity 0.3s ease,transform 0.3s ease}}
.sticky-title .dot{{color:var(--accent);margin:0 0.3em}}
.tabnav.stuck .sticky-title{{opacity:1;transform:translateY(0) scale(1);pointer-events:auto}}
.sticky-title:hover{{color:var(--accent)}}
.tabnav button{{font-family:'JBM', monospace;font-size:0.78rem;font-weight:700;letter-spacing:0.03em;text-transform:uppercase;background:none;border:none;color:var(--ink-faint);padding:0.7rem 1rem;cursor:pointer;border-bottom:2px solid transparent;margin-bottom:-1px;transition:color 0.15s ease,border-color 0.15s ease}}
.tabnav button:hover{{color:var(--ink)}}
.tabnav button.active{{color:var(--accent);border-bottom-color:var(--accent)}}
.tabnav button .n{{color:var(--ink-faint);font-weight:400;margin-left:0.35em}}
.panel{{display:none}}
.panel.active{{display:block}}
.scorecard{{display:grid;grid-template-columns:repeat(auto-fit,minmax(148px,1fr));gap:1px;background:var(--line);border:1px solid var(--line);margin-bottom:2.5rem;border-radius:2px;overflow:hidden}}
.score-cell{{background:var(--paper-raised);padding:0.85rem 0.9rem;cursor:pointer;transition:background 0.15s ease}}
.score-cell:hover{{background:var(--accent-soft)}}
.score-num{{font-family:'JBM', monospace;font-size:0.68rem;letter-spacing:0.08em;color:var(--ink-faint);font-variant-numeric:tabular-nums;display:flex;justify-content:space-between}}
.score-label{{font-size:0.8rem;font-weight:600;font-family:'JBM', monospace;margin:0.35rem 0 0.55rem;line-height:1.35;color:var(--ink)}}
.chip{{display:inline-flex;align-items:center;gap:0.35em;font-family:'JBM', monospace;font-size:0.64rem;font-weight:700;letter-spacing:0.05em;text-transform:uppercase;padding:0.2em 0.55em;border-radius:2px;white-space:nowrap}}
.chip::before{{content:'';width:6px;height:6px;border-radius:50%;display:inline-block}}
.chip-fixed{{background:var(--accent-soft);color:var(--accent-ink)}}
.chip-fixed::before{{background:var(--accent)}}
.chip-open{{background:var(--open-soft);color:var(--open-ink)}}
.chip-open::before{{background:var(--open)}}
.chip-pass{{background:var(--accent-soft);color:var(--accent-ink)}}
.chip-pass::before{{background:var(--accent)}}
.layout{{display:grid;grid-template-columns:210px 1fr;gap:3rem;align-items:start}}
.index{{position:sticky;top:3.5rem;z-index:1;display:flex;flex-direction:column;gap:0.15rem}}
.index-label{{font-family:'JBM', monospace;font-size:0.68rem;letter-spacing:0.1em;text-transform:uppercase;color:var(--ink-faint);margin-bottom:0.6rem}}
.index a{{display:block;font-family:'JBM', monospace;font-size:0.75rem;color:var(--ink-soft);text-decoration:none;padding:0.5rem 0 0.5rem 0.75rem;border-left:2px solid var(--line);line-height:1.5;transition:border-color 0.15s ease,color 0.15s ease}}
.index a .p{{display:block;font-size:0.62rem;color:var(--ink-faint);text-transform:uppercase;letter-spacing:0.06em;margin-top:0.15rem}}
.index a:hover{{color:var(--ink);border-left-color:var(--line-strong)}}
.index a.active{{color:var(--accent);border-left-color:var(--accent)}}
.posts{{min-width:0}}
article{{padding:2.5rem 0;border-bottom:1px solid var(--line);scroll-margin-top:1.5rem}}
article:first-child{{padding-top:0}}
article:last-of-type{{border-bottom:none}}
.eyebrow{{display:flex;align-items:center;gap:0.75rem;flex-wrap:wrap;font-family:'JBM', monospace;font-size:0.72rem;letter-spacing:0.06em;text-transform:uppercase;color:var(--ink-faint);margin-bottom:0.75rem;font-variant-numeric:tabular-nums}}
.eyebrow .num{{color:var(--accent);font-weight:700}}
article h2{{font-size:clamp(1.3rem,2.6vw,1.65rem);margin-bottom:1.1rem;max-width:28ch}}
article p{{max-width:66ch;margin:0 0 1.15rem;font-size:1.03rem;color:var(--ink)}}
article p strong{{font-weight:600}}
.callout{{max-width:66ch;border-left:3px solid var(--accent);border-radius:0;background:var(--accent-soft);padding:0.9rem 1.1rem;margin:0 0 1.15rem;font-size:0.97rem;color:var(--accent-ink)}}
code{{font-family:'JBM', ui-monospace, monospace;font-size:0.86em;background:var(--code-bg);padding:0.1em 0.35em;border-radius:2px;color:var(--ink);word-break:break-word}}
footer{{margin-top:2.5rem;padding-top:1.5rem;border-top:1px solid var(--line-strong);font-family:'JBM', monospace;font-size:0.72rem;color:var(--ink-faint);letter-spacing:0.03em}}
.section-intro{{max-width:70ch;color:var(--ink-soft);font-size:0.98rem;margin-bottom:2rem}}
.card-grid{{display:grid;grid-template-columns:1fr;gap:1px;background:var(--line);border:1px solid var(--line);border-radius:2px;overflow:hidden}}
.item-card{{background:var(--paper-raised);padding:1.5rem 1.75rem;scroll-margin-top:1.5rem}}
.item-eyebrow{{display:flex;align-items:center;gap:0.65rem;flex-wrap:wrap;font-family:'JBM', monospace;font-size:0.68rem;letter-spacing:0.05em;text-transform:uppercase;color:var(--ink-faint);margin-bottom:0.6rem;font-variant-numeric:tabular-nums}}
.item-eyebrow .num{{color:var(--accent);font-weight:700}}
.area-tag{{background:var(--code-bg);color:var(--ink-soft);padding:0.15em 0.5em;border-radius:2px}}
.filed-badge{{padding:0.15em 0.5em;border-radius:2px;font-weight:700}}
.filed-yes{{background:var(--accent-soft);color:var(--accent-ink)}}
.filed-no{{background:var(--code-bg);color:var(--ink-faint)}}
.filed-ref{{margin-bottom:0.9rem}}
.filed-ref code{{font-size:0.72rem}}
.item-card h3{{font-size:1.08rem;margin-bottom:0.85rem;max-width:60ch}}
.item-card p{{max-width:66ch;margin:0 0 0.9rem;font-size:0.96rem;color:var(--ink)}}
.item-card .callout{{font-size:0.9rem}}
@media (max-width:820px){{
  .layout{{grid-template-columns:1fr}}
  .index{{position:static;flex-direction:row;overflow-x:auto;gap:0;border-bottom:1px solid var(--line);padding-bottom:0.25rem;margin-bottom:0.5rem}}
  .index a{{border-left:none;border-bottom:2px solid var(--line);padding:0.5rem 0.85rem;white-space:nowrap}}
  .index a .p{{display:none}}
  .index a.active{{border-left:none;border-bottom-color:var(--accent)}}
  .tabnav{{overflow-x:auto}}
}}
@media (prefers-reduced-motion:reduce){{.score-cell,.index a,.sticky-title,.tabnav{{transition:none}}}}
.bloglog :focus-visible{{outline:2px solid var(--accent);outline-offset:2px}}
</style>
<h2 class="sr-only" style="position:absolute;width:1px;height:1px;overflow:hidden;clip:rect(0 0 0 0);">Audos SDK build log, bug tracker, and feature request list.</h2>
<div class="bloglog">
<div class="wrap">
  <header class="masthead">
    <div>
      <div class="masthead-title">Audos SDK <span class="dot">&middot;</span> Build Log</div>
      <div class="masthead-sub">One continuous record across every product we've built on Audos. We surface what we find, good and bad, to the cohort and to Audos as we go.</div>
    </div>
    <div class="masthead-meta">
      {blog_count} entries &middot; {date_range}<br>
      {bugs_open} bugs open, {bugs_fixed} fixed<br>
      {fr_count} feature requests &middot; {exp_count} experiments &middot; {sdk_count} SDK docs
    </div>
  </header>
  <nav class="tabnav" id="tabNav">
    <span class="sticky-title" id="stickyTitle" title="Back to top">Audos SDK<span class="dot">&middot;</span>Build Log</span>
    <button data-panel="panel-blog" class="active">Blog<span class="n">{blog_count}</span></button>
    <button data-panel="panel-experiments">Experiments<span class="n">{exp_count}</span></button>
    <button data-panel="panel-bugs">Bugs<span class="n">{bugs_count}</span></button>
    <button data-panel="panel-fr">Feature Requests<span class="n">{fr_count}</span></button>
    <button data-panel="panel-sdk">SDK<span class="n">{sdk_count}</span></button>
  </nav>

  <section class="panel active" id="panel-blog">
    <div class="scorecard">{blog_scorecard}
    </div>
    <div class="layout">
      <nav class="index" id="indexNav">
        <div class="index-label">Entries</div>
{blog_index_nav}
      </nav>
      <main class="posts">{blog_articles}
      </main>
    </div>
  </section>

  <section class="panel" id="panel-experiments">
    <div class="section-intro">Deliberate, hypothesis-driven tests of a specific platform capability &mdash; distinct from a bug (something broken) or a feature request (something missing). Each one states what we expected, what we actually did, and what came back, including the times we got it wrong the first pass and had to re-test.</div>
    <div class="card-grid">{exp_cards}
    </div>
  </section>

  <section class="panel" id="panel-bugs">
    <div class="section-intro">Pulled from <code>BACKLOG.md</code> and <code>docs/platform/bug-reports/</code> &mdash; every bug found during real builds on Audos, one card per issue, whether we fixed it ourselves or it's still open.</div>
    <div class="card-grid">{bugs_cards}
    </div>
  </section>

  <section class="panel" id="panel-fr">
    <div class="section-intro">Things that don't exist yet on the platform, needed for a credible developer/power-user channel &mdash; distinct from bugs (nothing here is "broken," it's unbuilt). See <code>docs/platform/20-power-user-wishlist.md</code> for the full narrative behind these.</div>
    <div class="card-grid">{fr_cards}
    </div>
  </section>

  <section class="panel" id="panel-sdk">
    <div class="section-intro">Not a log &mdash; an index. The durable, generic platform knowledge this project has accumulated along the way, distilled out of the bugs, experiments, and posts above and kept as its own reference. Lives in <code>docs/platform/</code> in the SDK repo; each card below is a short pointer to what's actually in there.</div>
    <div class="card-grid">{sdk_cards}
    </div>
  </section>

  <footer>
    Full technical detail, verified capability matrix, and the platform issue backlog live in the Audos SDK repo (<code>docs/platform/</code>, <code>BACKLOG.md</code>). This page is generated from <code>blog/*.md</code>, <code>blog/bugs/*.md</code>, and <code>blog/feature-requests/*.md</code> by <code>scripts/build_blog.py</code> &mdash; edit the Markdown, not this file.
  </footer>
</div>
</div>
<script>
(function(){{
  var tabs = Array.prototype.slice.call(document.querySelectorAll('#tabNav button'));
  var panels = tabs.map(function(b){{ return document.getElementById(b.dataset.panel); }});
  tabs.forEach(function(btn, i){{
    btn.addEventListener('click', function(){{
      tabs.forEach(function(b){{ b.classList.remove('active'); }});
      panels.forEach(function(p){{ p.classList.remove('active'); }});
      btn.classList.add('active');
      panels[i].classList.add('active');
    }});
  }});

  var links = Array.prototype.slice.call(document.querySelectorAll('#indexNav a'));
  var sections = links.map(function(a){{ return document.getElementById(a.dataset.target); }}).filter(Boolean);
  function onScroll(){{
    var pos = window.scrollY + 120;
    var current = sections[0];
    sections.forEach(function(s){{ if (s.offsetTop <= pos) current = s; }});
    links.forEach(function(a){{ a.classList.toggle('active', current && a.dataset.target === current.id); }});
  }}
  window.addEventListener('scroll', onScroll, {{ passive: true }});
  onScroll();

  var tabNav = document.getElementById('tabNav');
  var masthead = document.querySelector('.masthead');
  function onStickyScroll(){{
    tabNav.classList.toggle('stuck', masthead.getBoundingClientRect().bottom <= 0);
  }}
  window.addEventListener('scroll', onStickyScroll, {{ passive: true }});
  onStickyScroll();

  var stickyTitle = document.getElementById('stickyTitle');
  if (stickyTitle) {{
    var reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches;
    stickyTitle.addEventListener('click', function(){{
      window.scrollTo({{ top: 0, behavior: reduceMotion ? 'auto' : 'smooth' }});
    }});
  }}
}})();
</script>
'''

if __name__ == "__main__":
    main()
