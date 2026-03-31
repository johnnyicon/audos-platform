# Installing the Audos Platform Skill

Install this skill into any project that is built on the Audos platform. The symlink points to this repo on your local machine, so it only needs to be set up once per developer environment.

## Install

```bash
mkdir -p /path/to/your-audos-project/.claude/skills
ln -s /Users/kanekoa/Workspace/audos-platform/skill /path/to/your-audos-project/.claude/skills/audos-platform
```

> The `.claude/skills/` folder and its symlinks are gitignored — they're local developer setup, not checked in.

## Add the Rule to CLAUDE.md

Add the following to the project's `CLAUDE.md`:

```markdown
## Audos Platform

This project runs on the Audos platform. When working on anything that touches:
- The `audos-workspace/` folder
- Any HTTP endpoint at `https://audos.com/api/hooks/execute/...`
- Database tables, server functions, apps, or platform integrations

→ Read and follow: `/Users/kanekoa/Workspace/audos-platform/skill/SKILL.md`
```

## Verify

```bash
ls /path/to/your-audos-project/.claude/skills/audos-platform/SKILL.md
```
