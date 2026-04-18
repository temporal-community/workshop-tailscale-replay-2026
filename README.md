# Temporal + Tailscale Workshop — Replay 2026

## Topology

```
╔══════════════════════════════════════════════════════════════════════╗
║                        Tailscale Tailnet                             ║
║                                                                      ║
║   temporal-dev:7233 (gRPC)       metrics-server:9100/metrics         ║
║   temporal-dev:8233 (UI)         node_exporter                       ║
║   [ Temporal dev server ]        [ Linux VM ]                        ║
║                                                                      ║
║   http://ai                                                          ║
║   [ Aperture — Claude API ]                                          ║
║                                                                      ║
║   ┌──────────────────────────────────────────────┐                   ║
║   │  lab-worker                                  │                   ║
║   │                                              │                   ║
║   │  • joins tailnet via tsnet (no Tailscale     │                   ║
║   │    install required on your machine)         │                   ║
║   │  • polls metrics-server every minute         │                   ║
║   │  • asks Aperture to analyze the metrics      │                   ║
║   │  • records result in Temporal workflow       │                   ║
║   └──────────────────────────────────────────────┘                   ║
╚══════════════════════════════════════════════════════════════════════╝
```

Everything except the worker is pre-provisioned by the workshop facilitators. You only run the worker.

## Prerequisites

- Go 1.26+
- A `TS_AUTHKEY` provided by your workshop facilitator

## Run the worker

```shell
git clone https://github.com/temporal-community/workshop-tailscale-replay-2026
cd workshop-tailscale-replay-2026

TS_AUTHKEY=tskey-auth-<provided-key> \
METRICS_URL=http://metrics-server:9100/metrics \
go run ./cmd/worker
```

The worker joins the tailnet as `lab-worker` and immediately starts a cron workflow that runs every minute.

## View results

Open the Temporal UI in your browser (must be connected to the tailnet):

```
http://temporal-dev:8233
```

Navigate to the `health-check` workflow. Each completed run shows a 2–3 sentence AI-generated health summary of the metrics server.

## Environment variables

| Variable        | Required | Default          | Description                          |
|-----------------|----------|------------------|--------------------------------------|
| `TS_AUTHKEY`    | yes      | —                | Tailscale auth key from facilitator  |
| `METRICS_URL`   | yes      | —                | node_exporter endpoint on the tailnet|
| `TEMPORAL_HOST` | no       | `temporal-dev:7233` | Temporal server (via tsnet)       |
| `AI_URL`        | no       | `http://ai`      | Aperture endpoint (via tsnet)        |
| `AI_MODEL`      | no       | `claude-haiku-4-5` | Claude model                       |

## Run tests (no tailnet needed)

```shell
go test ./...
```

---

## Local testing (Mac — all services on one machine)

Runs the full stack locally on a Mac that has Tailscale installed. Three terminal windows.

### Prerequisites

```shell
# Temporal CLI
brew install temporal

# temporal-ts-net extension (puts Temporal on the tailnet)
# Using this fork until the fix is merged into temporal-community/temporal-ts-net
git clone https://github.com/kartikb-tailscale/temporal-ts-net
cd temporal-ts-net
go install ./cmd/temporal-ts_net
cd ..

# Verify the extension is found
temporal help --all | grep ts-net

# node_exporter
brew install node_exporter
```

You need two auth keys — one for the `temporal-ts-net` extension (registers as `temporal-dev`) and one for the worker (registers as `lab-worker`). Generate ephemeral keys in the [Tailscale admin console](https://login.tailscale.com/admin/settings/keys).

### Terminal 1 — Temporal dev server

```shell
TS_AUTHKEY=tskey-auth-<key-1> temporal ts-net
```

On first run it joins the tailnet as `temporal-dev`. You'll see:

```
Tailnet gRPC: temporal-dev:7233
Tailnet UI:   http://temporal-dev:8233
```

### Terminal 2 — node_exporter

node_exporter binds to `localhost` by default. Pass `--web.listen-address` so it listens on the Tailscale interface:

```shell
node_exporter --web.listen-address=0.0.0.0:9100
```

Find your Mac's Tailscale hostname:

```shell
tailscale status | head -3
```

Verify it's reachable from the tailnet:

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

Open `http://temporal-dev:8233` in a browser. The `health-check` workflow starts immediately and runs every minute. Click into a completed run — the workflow result is the AI-generated health summary.
