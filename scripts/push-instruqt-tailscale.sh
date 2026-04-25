#!/bin/zsh
set -euo pipefail

REPO_ROOT="$(git -C "$(dirname "$0")" rev-parse --show-toplevel)"
INSTRUQT_DIR="$REPO_ROOT/instruqt"

# Patch track.yml for tailscale org
sed -i '' 's/^owner: temporal/owner: tailscale/' "$INSTRUQT_DIR/track.yml"
sed -i '' '/^id: /d' "$INSTRUQT_DIR/track.yml"
sed -i '' '/mason\.egger@temporal\.io/d' "$INSTRUQT_DIR/track.yml"

# Clear challenge IDs
for f in "$INSTRUQT_DIR"/*/assignment.md; do
    sed -i '' '/^id: /d' "$f"
done

# Push (instruqt requires cwd to be the track directory)
(builtin cd "$INSTRUQT_DIR" && instruqt track push --force)

# Restore so we don't accidentally commit the patched files
git -C "$REPO_ROOT" checkout -- instruqt/
