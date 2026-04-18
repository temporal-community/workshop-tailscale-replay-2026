# Temporal + Tailscale Workshop — Replay 2026

## Topology

```
╔══════════════════════════════════════════════════════════════════════╗
║                         Tailscale Tailnet                            ║
║                                                                      ║
║  ┌─────────────────────────┐   ┌──────────────────────────────────┐  ║
║  │  temporal-dev           │   │  metrics-server                  │  ║
║  │  Linux VM               │   │  Linux VM                        │  ║
║  │                         │   │                                  │  ║
║  │  temporal-ts-net        │   │  node_exporter (Docker)          │  ║
║  │  joins via tsnet        │   │  :9100/metrics                   │  ║
║  │  :7233 gRPC             │   │                                  │  ║
║  │  :8233 UI               │   │  Tailscale (system client)       │  ║
║  └─────────────────────────┘   │  joins as metrics-server         │  ║
║                                └──────────────────────────────────┘  ║
║  ┌─────────────────────────┐                                         ║
║  │  http://ai              │                                         ║
║  │  Aperture (Claude API)  │                                         ║
║  └─────────────────────────┘                                         ║
║                                                                      ║
║  ┌─────────────────────────────────────────────────────────────┐     ║
║  │  lab-worker  ←  YOU RUN THIS                                │     ║
║  │  joins tailnet via tsnet (no Tailscale install required)    │     ║
║  │  fetches metrics-server:9100 → asks Aperture → records in   │     ║
║  │  Temporal workflow every minute                             │     ║
║  └─────────────────────────────────────────────────────────────┘     ║
╚══════════════════════════════════════════════════════════════════════╝
```

The `temporal-dev` and `lab-worker` components join the tailnet via
[tsnet](https://pkg.go.dev/tailscale.com/tsnet) — no Tailscale install
required, just a `TS_AUTHKEY`. The `metrics-server` VM runs the system
Tailscale client alongside node_exporter in Docker.

---

## Workshop: Running the worker

Everything except the worker is pre-provisioned by the workshop
facilitators. You only need Go and an auth key.

### Prerequisites

- Go 1.26+
- `TS_AUTHKEY` provided by your workshop facilitator

### Run

```shell
git clone https://github.com/temporal-community/workshop-tailscale-replay-2026
cd workshop-tailscale-replay-2026

TS_AUTHKEY=tskey-auth-<provided-key> \
METRICS_URL=http://metrics-server:9100/metrics \
go run ./cmd/worker
```

The worker joins the tailnet as `lab-worker` and starts a cron workflow
that runs every minute.

### View results

Open the Temporal UI in your browser (you must be on the tailnet):

```
http://temporal-dev:8233
```

Navigate to the `health-check` workflow. Each completed run shows a
system info block and an AI-generated health summary of the metrics
server.

### Environment variables

| Variable        | Required | Default               | Description                        |
|-----------------|----------|-----------------------|------------------------------------|
| `TS_AUTHKEY`    | yes*     | —                     | Tailscale auth key (required on first run; reuses stored state after) |
| `METRICS_URL`   | yes      | —                     | node_exporter endpoint on tailnet  |
| `TEMPORAL_HOST` | no       | `temporal-dev:7233`   | Temporal server address            |
| `AI_URL`        | no       | `http://ai`           | Aperture endpoint                  |
| `AI_MODEL`      | no       | `claude-haiku-4-5`    | Claude model                       |

### Run tests (no tailnet needed)

```shell
go test ./...
```

---

## Infrastructure setup (facilitators)

### temporal-dev VM

This VM runs `temporal server start-dev` and exposes it on the tailnet
as `temporal-dev` via [temporal-ts-net](https://github.com/kartikb-tailscale/temporal-ts-net).

**Install Go 1.26+**

```shell
wget https://go.dev/dl/go1.26.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.26.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc && source ~/.bashrc
```

**Install the Temporal CLI**

```shell
curl -sSf https://temporal.download/cli.sh | sh
echo 'export PATH=$PATH:$HOME/.temporalio/bin' >> ~/.bashrc && source ~/.bashrc
```

**Build the temporal-ts-net extension**

```shell
git clone https://github.com/kartikb-tailscale/temporal-ts-net
cd temporal-ts-net
go install ./cmd/temporal-ts_net
cd ..
```

Verify the extension is found:

```shell
temporal help --all | grep ts-net
```

**Run**

```shell
TS_AUTHKEY=tskey-auth-<key> temporal ts-net
```

On first run it joins the tailnet as `temporal-dev`. You'll see:

```
Tailnet gRPC: temporal-dev:7233
Tailnet UI:   http://temporal-dev:8233
```

---

### metrics-server VM

This VM runs node_exporter in Docker and joins the tailnet as
`metrics-server` via the system Tailscale client.

**Install Tailscale**

```shell
curl -fsSL https://tailscale.com/install.sh | sh
tailscale up --authkey=tskey-auth-<key> --hostname=metrics-server
```

Verify it's on the tailnet:

```shell
tailscale status | grep metrics-server
```

**Install Docker**

```shell
curl -fsSL https://get.docker.com | sh
```

**Run node_exporter**

```shell
docker run -d \
  --name node-exporter \
  --restart unless-stopped \
  --net=host \
  --pid=host \
  -v /:/host:ro,rslave \
  prom/node-exporter \
  --path.rootfs=/host \
  --web.listen-address=0.0.0.0:9100
```

`--net=host` ensures node_exporter listens on the Tailscale interface so
other tailnet nodes can reach it.

Verify metrics are reachable on the tailnet:

```shell
curl http://metrics-server:9100/metrics | head -5
```

---

## Local dev testing (Mac — all services on one machine)

Runs the full stack locally on a Mac that has Tailscale installed.
Three terminal windows. You need two auth keys — one for `temporal-dev`,
one for `lab-worker`. Generate ephemeral keys in the
[Tailscale admin console](https://login.tailscale.com/admin/settings/keys).

### Prerequisites

```shell
# Temporal CLI
brew install temporal

# temporal-ts-net extension
git clone https://github.com/kartikb-tailscale/temporal-ts-net
cd temporal-ts-net
go install ./cmd/temporal-ts_net
cd ..

# Verify
temporal help --all | grep ts-net

# node_exporter
brew install node_exporter
```

### Terminal 1 — temporal-dev

```shell
TS_AUTHKEY=tskey-auth-<key-1> temporal ts-net
```

### Terminal 2 — node_exporter

node_exporter binds to `localhost` by default. Pass
`--web.listen-address` so it's reachable on the Tailscale interface:

```shell
node_exporter --web.listen-address=0.0.0.0:9100
```

Find your Mac's Tailscale hostname:

```shell
tailscale status | head -3
```

Verify it's reachable:

```shell
curl http://<your-mac-tailscale-hostname>:9100/metrics | head -5
```

### Terminal 3 — worker

```shell
TS_AUTHKEY=tskey-auth-<key-2> \
METRICS_URL=http://<your-mac-tailscale-hostname>:9100/metrics \
go run ./cmd/worker
```

### Verify

Open `http://temporal-dev:8233` in a browser. The `health-check`
workflow runs every minute. Click into a completed run to see the
AI-generated health summary.
