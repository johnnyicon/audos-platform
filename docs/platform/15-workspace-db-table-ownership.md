# Audos Workspace DB — Table Ownership & Cleanup Analysis

**Date**: 2026-04-17  
**Purpose**: Classify all 20 tables as (a) deprecated Audos-app tables we can request cleanup on, (b) Audos platform tables to keep, or (c) ambiguous — needs Otto confirmation before deciding.

The daemon now handles all Throughline-specific workflow. Audos built several apps (Briefing, Guest Prep, Studio, Podcast Setup) that served those workflows before the daemon existed. Those app tables are now stale.

---

## Classification

### DEPRECATED — Audos-built app tables, replaced by the daemon

All of these belong to specific Audos apps that are no longer in use. The daemon handles everything they did. Safe to request removal.

| Table | Rows | What it was for | What replaced it |
|---|---|---|---|
| `app_briefing_podcast_profiles` | 0 | Briefing app — podcast show config | `podcast_config` in daemon |
| `app_briefing_research_sessions` | 0 | Briefing app — guest research + RoS | `episode_assets` (research type) in daemon |
| `app_briefing_ros_versions` | 0 | Briefing app — run-of-show version history | `episode_assets` (arc type) in daemon |
| `app_guest_prep_podcast_profiles` | 1 | Guest Prep app — duplicate podcast config | Same as above (1 stale row — old SG2GG config) |
| `app_guest_prep_research_sessions` | 0 | Guest Prep app — research per guest | `episode_assets` (research) + `episode_sources` in daemon |
| `app_guest_prep_ros_versions` | 0 | Guest Prep app — run-of-show version history | Arc in daemon |
| `app_podcast_setup_profiles` | 0 | Onboarding wizard — show setup | Settings page in Throughline + daemon `podcast_config` |
| `app_studio_episodes` | 0 | Studio app — episode tracking | `episodes` table in daemon |
| `app_studio_content` | 0 | Studio app — social content per episode | Future Throughline content pipeline |
| `app_studio_generated_content` | 0 | Studio app — generated copy | Same |
| `app_studio_time_tracking` | 0 | Studio app — time-saved metrics | Not replaced (was vanity metric) |

**Total deprecated**: 11 tables, 1 data row (1 stale podcast config from Guest Prep app).

---

### KEEP — Audos platform features, still active or depended on

These belong to Audos's platform layer, not to specific Throughline workflow apps. Removing them would break Audos functionality.

| Table | Rows | Why keep |
|---|---|---|
| `app_speakers` | 3 | Audos's speaker/host registry. John Gonzales (host) + SG2GG (brand) + Jess Thorne (guest). Voice profile system and caption generation reference this. |
| `app_voice_profiles` | 2 | Audos voice training infrastructure. John + SG2GG brand profiles. Throughline Signature depends on this — training happens through the Audos UI and is stored here. |
| `app_voice_refinements` | 0 | Stores refinement history when voice profiles are corrected via conversational feedback. Part of the voice platform even if currently empty. |
| `app_outreach_leads` | 11 | Audos's outreach CRM — generates and manages sales leads. This is Audos running their own feature for Kane, not a Throughline workflow table. |
| `app_dashboard_activity` | 2 | Audos's activity feed — tracks API calls, events in the workspace. Platform-level logging. |

**Total to keep**: 5 tables.

---

### AMBIGUOUS — Needs Otto to classify before deciding

These could belong to either the deprecated app layer OR the Audos platform depending on how Audos architected them.

| Table | Rows | Why ambiguous |
|---|---|---|
| `app_reels` | 1 | One draft reel (Jess Thorne, "Why DonorsChoose Exists"). Could be: (a) a deprecated Reels app that was built for Throughline, or (b) Audos's Reels platform feature for social clips. If (b), removing it breaks Audos's reel pipeline. |
| `app_reel_captions` | 0 | Caption content attached to reels. Same ambiguity as `app_reels`. If Audos's caption product reads from this, keep. |
| `app_generated_captions` | 0 | General-purpose caption generation table. Could be used by the Audos captions feature (platform) or was only used by deprecated studio/reel workflows. |
| `app_linked_references` | 2 | Stores fetched URL content (both rows are fetches of trythroughline.com from 2026-03-31). If this was the Briefing/Guest Prep app's research cache, it's deprecated. If it's a shared platform cache, keep. |

**Total ambiguous**: 4 tables.

---

## Summary

| Category | Count | Tables |
|---|---|---|
| Deprecated (request removal) | 11 | briefing_*, guest_prep_*, podcast_setup, studio_* |
| Keep (Audos platform) | 5 | speakers, voice_profiles, voice_refinements, outreach_leads, dashboard_activity |
| Confirm with Otto | 4 | reels, reel_captions, generated_captions, linked_references |

---

## Questions for Otto (append to the existing Otto prompt)

1. **The 11 deprecated app tables** (Briefing, Guest Prep, Studio, Podcast Setup) — these were created by Audos apps that are no longer in use. The daemon now handles all of that workflow. Can we drop these tables from the workspace schema, or does Audos need to handle the cleanup? If Audos needs to handle it, can we request that?

2. **`app_reels` / `app_reel_captions` / `app_generated_captions`** — are these part of an Audos platform feature (the Reels/Captions product) that Audos's own systems write to, or were they created exclusively for a specific app that's now deprecated?

3. **`app_linked_references`** — is this a shared platform cache used by multiple Audos features, or was it only used by the Briefing/Guest Prep apps for research URL fetching?

4. **If we can drop tables ourselves**: what's the safe process? Can we run `DROP TABLE` directly, or is there a migration/API path that Audos prefers?

---

## What this means for Atlas

Atlas manages the **daemon's Postgres** (Throughline's own DB), not the Audos workspace DB. These are two separate databases. Atlas does not need to know about or manage anything in the Audos schema.

The cleanup action here is not an Atlas migration — it's either:
- A direct `DROP TABLE` in the Audos workspace schema (if Otto confirms we have that right), or
- A support request to the Audos team to clean up deprecated app tables.

Once the Audos workspace is clean (platform tables only), the two DBs have clear ownership:
- **Audos workspace DB**: voice profiles, speakers, outreach, activity — Audos-managed.
- **Daemon DB (throughline_*)**: episodes, contacts, communications, arc, research, sources — Throughline-owned, Atlas-managed.
