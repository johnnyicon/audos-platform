# Genesis Space & the Custom-UI Ceiling

*Findings verified 2026-07-12 by having Otto inspect a live workspace ("DoKnow", `workspace-156396`).
Where the platform source wasn't readable from the workspace, that's called out explicitly.*

## Terminology: workspace vs. app (read this first)

Two distinct levels, and it's easy to blur them in conversation — worth being precise, because the fix
that applies at one level does not automatically apply at the other:

- **Workspace** — the top-level container. One per business/project (e.g. `workspace-156396`, the
  DoKnow workspace). Holds the shared shell (`Desktop.tsx`, `EmailGate.tsx`), shared theme tokens, the
  shared database, and shared config (`config.json`). There has been exactly **one** DoKnow workspace
  for the entire duration of this build — nothing in this SDK's work ever created a second one.
- **App** — a smaller unit living *inside* a workspace (e.g. `doknow-app`, `doknow-home`,
  `one-good-thing`; each is `apps/<Name>/App.tsx`). A workspace can hold many apps. Deleting or creating
  an app never touches the workspace it lives in.

This matters because of an easy trap: **a fix made once to the shared shell (`Desktop.tsx`) applies to
every app in that workspace, permanently — but a fix made to one app's own component file applies to
that app only**, and has to be re-briefed for the next one. See
`22-eliminating-the-chat-shell-playbook.md` ("Migrating an EXISTING app in place") for the concrete case
that surfaced this distinction, and `blog/0011-taking-the-red-pill.md` for the narrative version.

## Why every Audos build looks the same

Outputs converge visually across unrelated businesses **by design**, for three compounding reasons:

1. **A shared genesis space.** Every workspace is cloned from one central "genesis space." Evidence:
   the served landing shell carries a version marker —
   `EMAIL_GATE_VERSION = 101 // Landing-page shell aligned with legacy landing generation structure` —
   a strong tell that the shell is a shared, centrally-versioned template propagated to all workspaces.
2. **Token-based theming.** Styling flows through a theme-token system (`DesktopThemeTokens`,
   consumed via `SpaceRuntimeContext`/`useSpaceRuntime`). Components read brand tokens instead of
   hard-coded values, so a re-skin changes paint, not structure.
3. **Agents are told to fit the shell.** The build agents are instructed to produce apps that **fit
   the genesis shell and inherit brand tokens**. That instruction is the direct cause of the sameness:
   a learning app and a plumbing business inherit the same skeleton and differ only in tokens + content.

Net: the similarity is a **feature** for Audos (consistency + one upgrade path for all workspaces) that
reads as a **bug** when you want a distinctive product UI.

## What the genesis shell actually renders (observed in-browser 2026-07-12)

Viewing a signed-in workspace directly, the default surface the genesis shell produces is:

- **A ChatGPT-style chat as the primary surface** — a left sidebar with "New conversation" + a
  **Conversations** list, and a center **chat thread** with an "Ask me anything… (paste or drop images
  here)" composer.
- **Feature apps open as right-docked side panels**, listed under an **Apps** section in the sidebar
  with "Open app" links. The chat assistant literally responds to a request ("build me a course…") by
  **telling the user to open the relevant app** rather than doing it inline.

So out of the box, an Audos product is **"a chatbot with app-launcher side panels"** — there is no
default home/shelf/dashboard surface. This is the concrete reason Audos builds "feel like GPT with apps
on the side." Escaping it (a focused, single full-screen app) is the custom-UI ceiling question below.
For the running capability/limitation matrix and re-verification discipline, see
`19-capabilities-and-limitations.md`.

## Fixed vs flexible

| Layer | Status | Notes |
|-------|--------|-------|
| Space runtime / hosting contract (`SpaceRuntimeContext`, how `/space/…` is served, persistence API) | **FIXED** | Platform-owned; cannot change. |
| Desktop shell (`Desktop.tsx`), landing shell (`EmailGate.tsx`), dock / window manager | **Cloned, shared, editable-but-risky** | These files live in *your* workspace and are editable (the rename tooling walks + rewrites them), but they're the centrally-versioned skeleton — custom edits can collide with future template upgrades. |
| Theme tokens (colors, logo, typography), landing copy/layout in EmailGate, the mini-apps themselves | **FLEXIBLE** | The well-paved path. "Style the cloned genesis space" = apply brand tokens + business copy to the inherited shell. |
| Build-agent system prompts & formal design-system spec | **FIXED / not inspectable** | Platform-internal; `search_platform_code` returned "no searchable directories" from the workspace. Otto could not (and did not) quote them. |

## The mini-apps convention

Features are built as separate per-app React components (`apps/<Name>/App.tsx`) wrapped in the shared
dock/window shell. This is a **strong default, not a hard platform lock**, because:

- each app is an independently-editable unit — the `delegate_app_edit` / Cursor path operates per app;
- the dock/window shell is shared plumbing that upgrades centrally;
- it maps cleanly onto "add a feature = add an app."

The shell files are workspace-owned and editable, which is what makes the next point possible.

## The custom-UI ceiling (higher than "colors and copy")

**You can build a single full-screen app with its own internal navigation — no dock, no floating
windows — i.e. a real product dashboard.** Otto: *"a legitimate target, not a fight against the
platform."* The app build path supports mounting an app at a route where it owns its full render tree.

Two routes to a chrome-free single-app experience:
1. Build the app to occupy the full viewport at its route (lighter; stays on the paved path).
2. Edit the shell (`Desktop.tsx` / `EmailGate.tsx`) so the space boots straight into that one app
   instead of the dock (heavier; touches the shared, centrally-versioned skeleton — see the risk above).

Caveat: the exact mount/routing behavior for the *fully* chrome-free variant should be validated
against the genesis source during a real build (that source wasn't readable from the workspace at
time of writing).

## Tech stack the build agents actually use (verified from served source)

- **React + TypeScript**, one component per app (`apps/<Name>/App.tsx`).
- Icons: **`lucide-react`**.
- Styling: **theme-token driven** (`DesktopThemeTokens`) — a primary driver of visual convergence.
- Persistence: standardized on the **`useWorkspaceDB`** pattern.

## Practical guidance

- For a distinctive product, don't fight the tokens — **escape the shell**: brief the build for a
  single full-screen app that owns its layout and internal nav, rather than N dock apps.
- Expect anything that replaces the shell/dock/nav to be **editable-but-shared** territory: doable,
  but less supported and exposed to template-upgrade collisions. Weigh that against building the real
  UI off-platform (see the DoKnow evaluation in `doknow-kb/decisions/0001-audos-as-builder-evaluation.md`).
- The build-agent prompts/design-system are not externally inspectable; don't document guesses as fact.
