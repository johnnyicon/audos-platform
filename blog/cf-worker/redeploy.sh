#!/bin/bash
set -euo pipefail
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$REPO_ROOT"
python3 scripts/build_blog.py
cp blog/dist/index.html blog/cf-worker/public/index.html
cd blog/cf-worker
npx --yes wrangler@4 deploy
