# Persistent infrastructure

The workshop needs one long-lived VPS that runs `temporal-ts-net` and joins the Tailscale network as `temporal-dev`. Attendees never provision anything here. This is for whoever is hosting the workshop.

## The VPS

### Provision

- **Provider** - DigitalOcean, Hetzner, Linode, or equivalent. Any provider that gives you a plain Linux VM will do.
- **Specs** - **2 vCPU, 4 GB RAM** is enough for ~200 concurrent worker connections in a workshop setting. Step up to 4 vCPU / 8 GB only if you expect >200 attendees or are running the server under other workloads.
- **OS** - Ubuntu 22.04+ or Debian 12+. The instructions below assume a Debian-family distro.

!!! warning "This is the dev server, not a production Temporal Service"
    `temporal server start-dev` runs as a single process with all four Temporal services embedded, backed by SQLite with `NumHistoryShards: 1` and a single SQLite writer. The [docs](https://docs.temporal.io/cli/server#start-dev) are explicit that it's **not intended for production**. For a workshop it's great (one binary, no external database), but don't point long-lived production traffic at it. For production, use [Temporal Cloud](https://docs.temporal.io/cloud) or a [self-hosted cluster](https://docs.temporal.io/self-hosted-guide) with MySQL, Postgres, or Cassandra.

### Why these numbers

- **gRPC connections are cheap.** A worker's long-poll streams consume a few KB each and almost no CPU when idle. 200 idle pollers fit comfortably in a small VM.
- **SQLite writes are the real ceiling.** SQLite allows one writer at a time, so sustained workflow-completion rate is capped by disk and SQLite serialization, not by cores. Burst load from a workshop exercise (everyone firing the agent loop at once) is usually fine because individual workflows are short-lived.
- **Memory is mostly for workflow history cache.** For the workshop's short, linear agent workflows, history stays small.

If you're building something that runs for longer than a workshop or exposes the server to arbitrary traffic, switch to a real deployment before you hit the limits.

### Install the Temporal CLI

```shell
curl -sSf https://temporal.download/cli.sh | sh
echo 'export PATH=$PATH:$HOME/.temporalio/bin' >> ~/.bashrc
source ~/.bashrc
temporal --version   # must be v1.6.0+ for extension support
```

### Install `temporal-ts-net`

A prebuilt binary for your architecture, no Go toolchain required:

```shell
curl -sSfL https://raw.githubusercontent.com/temporal-community/temporal-ts-net/main/install.sh | sh
temporal help --all | grep ts-net   # verify the extension is found
```

The installer picks the right binary for `amd64` or `arm64`.

### Get a Tailscale auth key for the server

In the Tailscale admin console, under **Settings > Keys**:

1. **Reusable** - No (this key is for one machine, the VPS)
2. **Pre-authorized** - Yes
3. **Tags** - `tag:temporal-dev` (or whatever your ACLs use for the server)
4. **Expiration** - 90 days is a reasonable default

Copy the key. It starts with `tskey-auth-`.

### Start the server

```shell
export TS_AUTHKEY="tskey-auth-..."

temporal ts-net \
    --db-filename /var/lib/temporal/workshop.db \
    --max-connections 500 \
    --connection-rate-limit 50
```

The server becomes reachable at:

- **`temporal-dev:7233`** - gRPC (workers and the CLI connect here)
- **`temporal-dev:8233`** - Web UI (attendees view workflows here)

`--db-filename` persists workflow data across restarts. Skip it and everything is in-memory: fine for rehearsal, not fine for a live workshop.

About those flags (they come from `temporal-ts-net`, not the dev server itself):

- **`--max-connections 500`** - caps simultaneous tsnet listener accepts. For ~200 attendees with 2-3 worker processes each, 500 is ample headroom. Bump it if you're running a larger session.
- **`--connection-rate-limit 50`** - caps new connections per second. Protects against boot-storm fan-in when the whole room starts their workers in the same 30 seconds. 50/s drains 200 attendees in under 4 seconds.

For bigger cohorts (say, 500+), raise these proportionally and consider stepping up to 4 vCPU / 8 GB so SQLite's single writer isn't the bottleneck during peak bursts.

### Running as a service

For a real workshop, run it as a systemd unit so it survives reboots and crashes:

```shell
sudo tee /etc/systemd/system/temporal-tsnet.service << 'EOF'
[Unit]
Description=Temporal Dev Server on Tailscale
After=network.target

[Service]
Type=simple
Environment=TS_AUTHKEY=tskey-auth-...
ExecStart=/usr/local/bin/temporal ts-net --db-filename /var/lib/temporal/workshop.db --max-connections 500 --connection-rate-limit 50
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable --now temporal-tsnet
sudo systemctl status temporal-tsnet
```

## The attendee auth key

Separate from the server's key, you'll want a **reusable** auth key for attendee machines:

1. **Reusable** - Yes
2. **Pre-authorized** - Yes
3. **Tags** - `tag:workshop`
4. **Expiration** - align with the workshop duration (for example, 1 day for a one-day event)

Distribute this key however fits the session: an Instruqt lifecycle script, a slide, or a QR code at the door. Anyone with the key can join the tailnet, so treat it like a workshop-session secret.

## Aperture

Aperture is Tailscale's API gateway. The workshop uses it to hold the shared OpenAI key and rate-limit per attendee identity. Provisioning Aperture is a Tailscale-side concern; see the [Aperture documentation](https://docs.tailscale.com/aperture) for the current setup process.

What the workshop repo expects:

| Variable | Value |
|---|---|
| Aperture hostname on the tailnet | `ai` (so `http://ai/v1` resolves) |
| OpenAI-compatible endpoint | `POST /v1/responses` |
| Auth model | Tailscale identity. Aperture reads the caller's tailnet identity; attendees send no `Authorization` header |
| Rate-limit policy | Per-identity (demo slide shows ~3/10 requests; tune for your audience) |

If Aperture isn't available, attendees can fall back to a direct OpenAI key by overriding the `openai_base_url` and `openai_api_key` env vars in their `temporal.toml` profile.

## ACLs

Keep the ACL policy tight enough that an attendee can't spoof the server. A minimal shape:

- `tag:temporal-dev` - the VPS. Accepts inbound on `:7233` and `:8233` from `tag:workshop`.
- `tag:aperture` - the Aperture host. Accepts inbound on `:80` (or your chosen port) from `tag:workshop`.
- `tag:workshop` - every attendee machine. Can reach the two tags above, can't reach each other.

Blocking attendee-to-attendee traffic prevents the "oops someone found a dev server on their classmate's laptop" problem during the workshop.

## Teardown

After the workshop:

1. Rotate or delete the attendee auth key in the Tailscale console.
2. Stop the systemd unit: `sudo systemctl disable --now temporal-tsnet`.
3. Decide whether to preserve the workflow history in `/var/lib/temporal/workshop.db` (useful for post-workshop debugging) or wipe it.
4. If the VPS was purpose-built, destroy it. There's no ongoing state worth keeping.
