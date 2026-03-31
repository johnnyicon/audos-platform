# Installing the Audos Platform Skill

## Global Install (recommended)

Symlinks the skill into Claude Code's global skills directory so it's available in any project:

```bash
ln -s /Users/kanekoa/Workspace/audos-platform/skill ~/.claude/skills/audos-platform
```

## Per-Project Install

Symlinks the skill into a specific project's Claude skills directory:

```bash
mkdir -p /path/to/your-project/.claude/skills
ln -s /Users/kanekoa/Workspace/audos-platform/skill /path/to/your-project/.claude/skills/audos-platform
```

## Add the Rule to CLAUDE.md

Add the following to the project's `CLAUDE.md` so Claude Code automatically knows to use this skill when working with Audos:

```markdown
## Audos Platform

This project runs on the Audos platform. When working on anything that touches:
- The `audos-workspace/` folder
- Any HTTP endpoint at `https://audos.com/api/hooks/execute/...`
- Database tables, server functions, apps, or platform integrations

→ Read and follow: `/Users/kanekoa/Workspace/audos-platform/skill/SKILL.md`
```

## Verify

After installing, confirm the symlink resolves:

```bash
ls ~/.claude/skills/audos-platform/SKILL.md
```
