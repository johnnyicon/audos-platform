---
date: 2026-07-12
product: DoKnow
status: open
label: Chat-shell diagnosis
---

# It felt like ChatGPT with apps on the side

Once we could actually log in and click around, the reaction was immediate and blunt: this doesn't
feel like an app. It's not a design complaint — the palette and type were fine. It's structural.

Open the signed-in workspace and the primary surface is a chat window: a conversation list on the
left, a message thread in the center, an "Ask me anything" composer at the bottom. Type "build me a
course on manifestation" and the Coach doesn't build anything — it replies with a numbered list
telling you to go open the Course Builder app yourself. The apps we'd built — Course Builder, Lesson
Player, Coach Queue — open as right-docked side panels the chatbot points you to. There's no home
screen, no shelf of courses, no visible streak or progress, no "here's your next lesson." You drive
everything by hand, which is close to the opposite of what DoKnow is supposed to do.

Digging into why: every Audos workspace is cloned from one shared "genesis space" — a version marker
in the landing shell's own source (`EMAIL_GATE_VERSION = 101 // aligned with legacy landing generation
structure`) gave it away. The shell, the dock, the chat — all shared, centrally versioned plumbing.
Build agents are explicitly steered to "fit the genesis shell and inherit brand tokens." That's a
sane default for consistency and one-shot platform upgrades. It's also exactly why a learning app and
a plumbing business come out of Audos looking like siblings, and why our own three apps read as
disconnected tools glued onto someone else's chatbot rather than one coherent product.

> There was a second, deeper problem underneath the chrome. Every course Course Builder generated came
> back with a title, a description, one tier, one module — and zero lessons. Not thin lessons: none.
> The Lesson Player for those courses just said "this course has no lessons yet." So even setting the
> shell aside, the actual product value — a leveled course you can learn from — wasn't being produced.

Two separate problems, and worth naming clearly: a shell that doesn't feel like a product, and a
generator that doesn't generate the product. The next two posts are us going after each of them, in
the order we could actually test them.
