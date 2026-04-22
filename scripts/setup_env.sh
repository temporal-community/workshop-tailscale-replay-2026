#!/usr/bin/env bash
# ABOUTME: Instruqt lifecycle script to set up the workshop environment.
# ABOUTME: Clones the repo, installs dependencies, configures temporal.toml and shell env vars.

set -euo pipefail

REPO_URL="https://github.com/temporal-community/workshop-tailscale-replay-2026.git"
WORKSHOP_DIR="/root/workshop"

echo "==> Setting up workshop environment..."

# Clone the repository
if [ ! -d "$WORKSHOP_DIR" ]; then
    git clone "$REPO_URL" "$WORKSHOP_DIR"
else
    cd "$WORKSHOP_DIR" && git pull
fi

cd "$WORKSHOP_DIR"

# Derive user ID from hostname (Instruqt gives each VM a unique hostname)
WORKSHOP_USER_ID=$(hostname | tr '[:upper:]' '[:lower:]')

# Export environment variables for all shells
cat >> ~/.bashrc << EOF

# Workshop environment variables
export TEMPORAL_PROFILE=tailnet
export WORKSHOP_USER_ID=${WORKSHOP_USER_ID}
export OPENAI_BASE_URL=http://ai/v1
export OPENAI_API_KEY=workshop-token
EOF

# Create Temporal environment config with the tailnet profile
mkdir -p ~/.config/temporalio
cat > ~/.config/temporalio/temporal.toml << 'EOF'
# Temporal Environment Configuration
# Created by workshop setup script

[profile.default]
address = "localhost:7233"
namespace = "default"

[profile.tailnet]
address = "temporal-dev:7233"
namespace = "default"
EOF

# Install Python dependencies
uv sync

echo "==> Workshop environment ready!"
echo "    User ID: ${WORKSHOP_USER_ID}"
echo "    Temporal profile: tailnet (temporal-dev:7233)"
echo "    Run 'uv run scripts/verify_setup.py' to verify connectivity."
