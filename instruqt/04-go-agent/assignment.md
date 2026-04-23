---
slug: go-agent
id: cngfnpse6nrs
type: challenge
title: Metrics Watcher
teaser: Run a Go metrics watcher that joins the tailnet via tsnet and asks
  Claude via Aperture for a health summary on a Temporal Schedule
notes:
- type: text
  contents: |-
    # Metrics Watcher

    The final workshop activity. A finished Go worker that joins the
    tailnet via `tsnet`, scrapes `node_exporter` metrics off a
    `metrics-server` node, and asks Claude (via Aperture) for a
    plain-English health summary on a Temporal Schedule. No code to
    edit this time — your job is to run it, watch the Schedule fire
    in the Temporal UI, and tune the cadence.
tabs:
- id: rtxkm5tpeory
  title: Code Editor
  type: code
  hostname: workshop
  path: /root/workshop/exercises/04_go_agent/practice
- id: ojtn3trxjmsm
  title: Worker
  type: terminal
  hostname: workshop
  workdir: /root/workshop
- id: 2dllsneptu75
  title: Starter
  type: terminal
  hostname: workshop
  workdir: /root/workshop
- id: vbmyarbcjzjs
  title: Temporal UI
  type: service
  hostname: workshop
  port: 8233
- id: c4xnvxeanmfp
  title: Aperture UI
  type: service
  hostname: workshop
  port: 80
difficulty: advanced
timelimit: 1800
enhanced_loading: null
---

# Exercise 4: Metrics Watcher

The same `tsnet` pattern from Exercise 2, but against real services: pull `node_exporter` metrics off another tailnet node, ask Claude (via Aperture) for a plain-English health summary, and run it all on a Temporal Schedule. The code is already complete — you'll run it, watch it on the Temporal UI, and tune the cadence.

> **Last activity of the workshop.** There are no TODOs to fill in here. This exercise is lighter on coding and heavier on observing the system in action — a walkthrough of how the pieces fit together in a real, scheduled, durable agent.

> **Not on the tailnet?** If you joined late or `tailscale status` shows **Logged out**, run this in the **Worker** terminal first:
>
> ```bash
> tailscale up --auth-key="$TS_AUTHKEY" --hostname="${WORKSHOP_USER_ID}-env"
> ```

## Topology

```
┌──────────────────┐  tailnet   ┌────────────────────┐
│  Your worker     │ ─────────► │ temporal-dev:7233  │
│  (tsnet inside   │            └────────────────────┘
│   a Go process)  │  tailnet   ┌────────────────────┐
│                  │ ─────────► │ metrics-server     │
│                  │            │ (node_exporter)    │
│                  │  tailnet   └────────────────────┘
│                  │            ┌────────────────────┐   ┌─────────────────┐
│                  │ ─────────► │ Aperture (http://ai│ ► │ Anthropic API   │
└──────────────────┘            │ shared LLM gateway) │  │ (shared key)    │
                                └────────────────────┘   └─────────────────┘
```

The worker joins the tailnet *itself* via `tsnet` (just like Exercise 2), but this time its HTTP client reaches a real `metrics-server` node and its LLM calls go through Aperture to Claude. All three destinations are on the workshop tailnet — no public endpoints.

## What's already built

Open the **Code Editor** tab. The practice directory has:

- `main.go` — joins the tailnet via `tsnet`, dials Temporal through a `tsnet.Dial` context dialer, registers a Temporal Schedule that fires `HealthCheckWorkflow` every `HEALTH_CHECK_INTERVAL`.
- `activities.go` — `FetchMetrics` scrapes `node_exporter`; `AnalyzeMetrics` asks Claude via the Anthropic SDK (pointed at Aperture) and returns a structured `HealthReport`.
- `workflow.go` — `HealthCheckWorkflow` chains the two activities.
- `activities_test.go`, `workflow_test.go` — offline tests using `httptest.Server`.

## Step 1: Move to the practice directory

In the **Worker** terminal:

```bash
cd exercises/04_go_agent/practice
go mod download
```

## Step 2: Run the worker

```bash
METRICS_URL=http://metrics-server:9100/metrics go run .
```

`WORKSHOP_USER_ID`, `TS_AUTHKEY`, and `APERTURE_URL` are already exported by the workshop setup, so you only need `METRICS_URL`.

First run takes 10–30 seconds while `tsnet` registers the node. You should see log lines like:

```
level=INFO msg="joined tailnet" hostname=<you>-ex4-metrics-worker-<5 random chars> userID=<you>
level=INFO msg="connected to temporal" host=temporal-dev:7233
level=INFO msg="metrics reachable" url=http://metrics-server:9100/metrics
level=INFO msg="created schedule" id=<you>-health-check-schedule interval=10m0s workflow=<you>-health-check
level=INFO msg="worker running" taskQueue=<you>-health-check
```

## Step 3: Watch the Schedule in the Temporal UI

Open the **Temporal UI** tab. Two places to look:

- **Schedules** → `<your-user-id>-health-check-schedule`. Shows the interval, the next fire time, and recent runs. The Schedule fires immediately on startup because the worker registers it with `TriggerImmediately`.
- **Workflows** → search for `<your-user-id>-health-check`. Each fired run has a completed row (its ID is suffixed with the schedule fire time) whose Result panel contains the `HealthReport` JSON.

## Step 4: Tune the cadence

`10m` is too slow to watch during the workshop. In the **Worker** terminal, `Ctrl+C`, then restart with a shorter interval:

```bash
HEALTH_CHECK_INTERVAL=2m METRICS_URL=http://metrics-server:9100/metrics go run .
```

Any Go duration works (`30s`, `5m`, `1h`). The worker deletes and recreates the Schedule on every startup, so changing the interval just means restarting.

## Step 5: Customize the Claude prompt (optional)

Open `activities.go` in the **Code Editor** tab, find `AnalyzeMetrics`. The prompt lives in a raw string. Change it — ask Claude to flag anything unusual, add a field to the `HealthReport` struct, whatever. Restart the worker and watch the next Schedule fire produce a different `HealthReport` in the UI.

## Step 6: Run the offline tests

The tests mock `node_exporter` and Aperture with `httptest.Server` — no tailnet needed:

```bash
go test ./...
```

## Environment variables

| Variable                | Required | Default               | Description                                                          |
|-------------------------|----------|-----------------------|----------------------------------------------------------------------|
| `WORKSHOP_USER_ID`      | yes      | (none)                | Prefixes hostname, task queue, schedule ID, and workflow ID.         |
| `TS_AUTHKEY`            | yes*     | (none)                | Tailscale auth key. Required on first run; `tsnet` reuses state after.|
| `METRICS_URL`           | yes      | (none)                | `node_exporter` endpoint on the tailnet.                             |
| `HEALTH_CHECK_INTERVAL` | no       | `10m`                 | Cadence as a Go duration (`30s`, `5m`, `1h`).                        |
| `TEMPORAL_HOST`         | no       | `temporal-dev:7233`   | Temporal server address.                                             |
| `APERTURE_URL`          | no       | `http://ai`           | Aperture endpoint; Anthropic SDK appends `/v1/messages` automatically.|
| `AI_MODEL`              | no       | `claude-haiku-4-5`    | Claude model.                                                        |

`*` = required on first run only; the `tsnet` state dir persists the node key.

## What you've learned

- `tsnet.Dial` works for both tailnet-internal HTTP (metrics) and gRPC (Temporal)
- Aperture is model-agnostic: the same gateway proxies Anthropic here and OpenAI in Exercise 3
- Temporal Schedules with `TriggerImmediately` fire once on creation, then every N, with the next fire visible in the UI
- All three backing services (Temporal, metrics, Aperture) are tailnet-only; Tailscale identity is the auth layer — no keys on your machine
