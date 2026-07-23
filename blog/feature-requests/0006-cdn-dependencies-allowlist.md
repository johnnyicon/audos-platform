---
date: 2026-07-13
priority: 3
status: not filed
label: Documented `cdnDependencies` allow-list
---

Neither we nor Otto itself could read the accepted `cdnDependencies` list from any documented source —
`search_platform_code` returned "no searchable directories" even when Otto ran it. We had to build a
disposable throwaway app that tried five libraries live just to learn that react-query, GSAP, Radix, and
three.js work, and react-three-fiber is flaky. A published allow-list would save that entire discovery
cycle for the next team.
