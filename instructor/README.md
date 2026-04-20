# Instructor Setup Guide

Instructions for setting up the workshop infrastructure before the event.

## VPS: Temporal Dev Server

### Provision

- Provider: DigitalOcean, Hetzner, or similar
- Specs: 4+ CPU, 8GB+ RAM (handles ~30-50 concurrent worker connections)
- OS: Ubuntu 22.04+ or Debian 12+

### Install Temporal CLI

```bash
curl -sSf https://temporal.download/cli.sh | sh
```

Verify: `temporal --version` (needs v1.6.0+ for extension support)

### Install temporal-ts-net

```bash
curl -sSfL https://raw.githubusercontent.com/temporal-community/temporal-ts-net/main/install.sh | sh
```

Verify: `temporal ts-net --help`

### Get a Tailscale Auth Key

From the Tailscale admin console for the workshop tailnet:

1. Go to Settings → Keys
2. Generate a new auth key:
   - Reusable: Yes
   - Pre-authorized: Yes
   - Tags: `tag:instructor` (or whatever ACL tag Kartik set up)
3. Copy the key

### Start the Server

Use the `run_server.sh` script in this directory, or run manually:

```bash
export TS_AUTHKEY="tskey-auth-..."

temporal ts-net \
    --db-filename /var/lib/temporal/workshop.db \
    --max-connections 2000 \
    --connection-rate-limit 200
```

The server will be accessible at:
- `temporal-dev:7233` (gRPC — workers and CLI connect here)
- `temporal-dev:8233` (Web UI — attendees view workflows here)

### Persistence

The `--db-filename` flag persists workflow data across restarts. Use this during the workshop for resilience. Without it, all data is in-memory and lost on restart.

### Running as a Service (Optional)

For the actual workshop, consider running via systemd:

```bash
sudo tee /etc/systemd/system/temporal-tsnet.service << 'EOF'
[Unit]
Description=Temporal Dev Server on Tailscale
After=network.target

[Service]
Type=simple
Environment=TS_AUTHKEY=tskey-auth-...
ExecStart=/usr/local/bin/temporal ts-net --db-filename /var/lib/temporal/workshop.db --max-connections 2000 --connection-rate-limit 200
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now temporal-tsnet
```

## Instruqt Track

### VM Requirements

Each attendee VM needs:
- Python 3.13 (`pyenv` or system package)
- `uv` (Python package manager)
- Go 1.22+ (for Exercise 4 stretch goal)
- Tailscale installed and connected to the workshop tailnet
- Web browser access (for Temporal UI at `http://temporal-dev:8233`)

### Lifecycle Script

The `scripts/setup_env.sh` script should run during the Instruqt VM provisioning:
1. Clones the workshop repo
2. Creates `.env` with per-VM values
3. Runs `uv sync` to install dependencies

### Tailscale Setup in Instruqt

Each VM needs a Tailscale auth key to join the workshop tailnet. Options:
1. **Pre-authorized reusable key** in the lifecycle script (simplest)
2. **Tailscale Connector** if Instruqt supports it
3. **Manual join** as a fallback (attendees run `tailscale up --authkey ...`)

Coordinate with Kartik on the auth key approach.

## Pre-Event Checklist

- [ ] VPS provisioned and temporal-ts-net running
- [ ] `temporal-dev:7233` reachable from a Tailscale-connected machine
- [ ] `temporal-dev:8233` (Web UI) accessible in a browser
- [ ] Aperture endpoint confirmed with Kartik
- [ ] Instruqt track tested end-to-end
- [ ] Workshop repo up to date
- [ ] Run `scripts/verify_setup.py` from an Instruqt VM — all checks pass

## Day-Of Checklist

- [ ] VPS is running, check: `tailscale ping temporal-dev`
- [ ] Temporal UI loads at `http://temporal-dev:8233`
- [ ] Aperture is healthy (check with Kartik)
- [ ] Instruqt track is published and invite links sent
- [ ] Run one test workflow from an Instruqt VM to verify end-to-end
- [ ] Slides loaded and ready

## Backup Plan

If the shared server fails mid-workshop:
- Attendees can run `temporal server start-dev` locally on their Instruqt VM
- Change `TEMPORAL_ADDRESS` in `.env` to `localhost:7233`
- The Aperture integration still works (LLM calls route through Tailscale)
- They lose the shared UI experience but can still complete exercises
