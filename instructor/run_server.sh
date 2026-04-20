#!/usr/bin/env bash
# ABOUTME: Starts the Temporal dev server on the Tailscale network via temporal-ts-net.
# ABOUTME: Run this on the VPS that hosts the shared workshop server.

set -euo pipefail

# Ensure TS_AUTHKEY is set
if [ -z "${TS_AUTHKEY:-}" ]; then
    echo "ERROR: TS_AUTHKEY environment variable is not set."
    echo "Get a reusable, pre-authorized auth key from the Tailscale admin console."
    echo ""
    echo "Usage:"
    echo "  export TS_AUTHKEY='tskey-auth-...'"
    echo "  ./run_server.sh"
    exit 1
fi

# Create data directory if needed
sudo mkdir -p /var/lib/temporal

echo "Starting Temporal dev server on the Tailscale network..."
echo "  gRPC:  temporal-dev:7233"
echo "  UI:    http://temporal-dev:8233"
echo ""

temporal ts-net \
    --db-filename /var/lib/temporal/workshop.db \
    --max-connections 2000 \
    --connection-rate-limit 200
