# Exercise 4 — deployment

Deployment runbook for future workshop creators: stand up the two VMs
that Exercise 4 depends on, run the worker against them, and a quick
Mac-only path for validating the whole setup in one place.

---

## Aperture (not self-hostable today)

Aperture is the only workshop component that isn't something you can
spin up yourself from the Tailscale side — at the moment you have to
sign up to get it provisioned into your tailnet. Sign-up lives at
[aperture.tailscale.com](https://aperture.tailscale.com/). Keep an eye
on that page for self-hosting support down the line.

Once Aperture is in your tailnet, the default `APERTURE_URL=http://ai`
works without any additional config; the worker resolves and reaches
it over the tailnet like any other node.

Everything else in this doc — `metrics-server`, `temporal-dev`, the
worker — is infrastructure you run yourself on whatever VM/host you
want.

---

## VM 1: `metrics-server`

Standard Linux VPS. Joins the tailnet via the system Tailscale client
and runs `node_exporter` via Docker Compose.

### 1. Install and start Tailscale

```shell
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up --authkey=tskey-auth-<key> --hostname=metrics-server
```

Verify:

```shell
tailscale status | grep metrics-server
```

### 2. Install Docker

```shell
curl -fsSL https://get.docker.com | sh
sudo systemctl enable --now docker
```

`enable --now` makes the daemon come up at boot. Combined with
`restart: unless-stopped` in the compose file, the container survives
reboots and crashes.

### 3. Run `node_exporter`

Drop in a `docker-compose.yml`:

```yaml
services:
  node-exporter:
    image: prom/node-exporter
    container_name: node-exporter
    restart: unless-stopped
    network_mode: host
    pid: host
    volumes:
      - /:/host:ro,rslave
    command:
      - --path.rootfs=/host
      - --web.listen-address=0.0.0.0:9100
```

Bring it up:

```shell
docker compose up -d
```

`network_mode: host` puts the listener on the Tailscale interface so
other tailnet nodes can reach it.

Verify from another tailnet node:

```shell
curl http://metrics-server:9100/metrics | head -5
```

---

## VM 2: `temporal-dev`

Standard Linux VPS. Runs `temporal ts-net` which embeds `tsnet` to
register the dev server as `temporal-dev` on the tailnet.

### 1. Install and start Tailscale

```shell
curl -fsSL https://tailscale.com/install.sh | sh
sudo tailscale up --authkey=tskey-auth-<key>
```

(The hostname here doesn't matter — `temporal-ts-net` registers its own
`temporal-dev` tailnet node separately.)

### 2. Install the Temporal CLI

```shell
curl -sSf https://temporal.download/cli.sh | sh
echo 'export PATH=$PATH:$HOME/.temporalio/bin' >> ~/.bashrc
source ~/.bashrc
```

### 3. Install the `temporal-ts-net` extension

```shell
curl -sSfL https://raw.githubusercontent.com/temporal-community/temporal-ts-net/main/install.sh | sh
temporal help --all | grep ts-net   # sanity check
```

### 4. Run the dev server

```shell
TS_AUTHKEY=tskey-auth-<key> temporal ts-net
```

First run takes a few seconds. Expect:

```
Tailnet gRPC: temporal-dev:7233
Tailnet UI:   http://temporal-dev:8233
```

For production-style persistence, run it under a systemd unit or a
process manager with the `TS_AUTHKEY` injected from a secret store.

---

## Run the worker

From any tailnet-connected machine (your laptop, an Instruqt VM, the
`metrics-server` VM itself — doesn't matter, the worker joins the
tailnet via `tsnet`):

```shell
cd exercises/04_go_agent/practice
go mod download

WORKSHOP_USER_ID=<your-handle> \
TS_AUTHKEY=tskey-auth-<key> \
METRICS_URL=http://metrics-server:9100/metrics \
go run .
```

Set `HEALTH_CHECK_INTERVAL=2m` (or any Go duration) for faster iteration.

Verify the Schedule fired: open `http://temporal-dev:8233`, look for
`<your-handle>-health-check-schedule`, and check for completed
`<your-handle>-health-check-<timestamp>` workflow runs with a
`HealthReport` in the result.

---

## Local-on-Mac quick validation

One Mac, three terminals, no remote VMs. Useful when you want to verify
all four pieces (`temporal-ts-net`, `node_exporter`, the worker, and
Aperture reachability) before you invest in a VPS setup.

**Prereqs**

- Tailscale installed and the Mac is on your personal tailnet.
- Aperture provisioned into that tailnet (see the section at the top
  of this doc) so `AnalyzeMetrics` can reach Claude.
- Two Tailscale auth keys ready — one for `temporal-ts-net`, one for
  the worker. Either generate them as **single-use** keys, or use
  reusable ones and revoke them in the admin console when you're done
  testing.

**Install**

```shell
brew install temporal node_exporter
curl -sSfL https://raw.githubusercontent.com/temporal-community/temporal-ts-net/main/install.sh | sh
temporal help --all | grep ts-net
```

**Terminal 1 — `temporal-dev`**

```shell
TS_AUTHKEY=tskey-auth-<key-1> temporal ts-net
```

**Terminal 2 — `node_exporter`**

```shell
node_exporter --web.listen-address=0.0.0.0:9100
```

`node_exporter` binds to `localhost` by default; the flag makes it
reachable on the tailnet interface. Grab the Mac's tailnet hostname:

```shell
tailscale status | head -1
```

Smoke-test:

```shell
curl http://<mac-tailnet-hostname>:9100/metrics | head -5
```

**Terminal 3 — worker**

```shell
cd exercises/04_go_agent/practice

WORKSHOP_USER_ID=<your-handle> \
TS_AUTHKEY=tskey-auth-<key-2> \
METRICS_URL=http://<mac-tailnet-hostname>:9100/metrics \
HEALTH_CHECK_INTERVAL=30s \
go run .
```

Open `http://temporal-dev:8233` and watch the Schedule fire. If the
first run produces a completed `HealthReport` in the UI, the entire
stack is wired correctly end-to-end.
