---
slug: go-agent
id: cngfnpse6nrs
type: challenge
title: Metrics Watcher
teaser: Run a Go metrics watcher that joins the tailnet via tsnet and asks Claude
  via Aperture for a health summary on a Temporal Schedule
notes:
- type: text
  contents: |-
    # Metrics Watcher

    The final workshop activity. A finished Go Worker that joins the
    `tailnet` via `tsnet`, scrapes `node_exporter` metrics off a
    `metrics-server` node, and asks Claude (via Aperture) for a
    plain-English health summary on a Temporal Schedule. No code to
    edit this time. Your job is to run it, watch the Schedule fire
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

The same `tsnet` pattern from Exercise 2, but against real services. The Worker pulls `node_exporter` metrics off another `tailnet` node, asks Claude (via Aperture) for a plain-English health summary, and runs it all on a Temporal Schedule. The code is already complete, you'll run it, watch it on the Temporal UI, and tune the cadence.

> **Last activity of the workshop.** There are no TODOs to fill in here. This exercise is lighter on coding and heavier on observing the system in action, a walkthrough of how the pieces fit together in a real, scheduled, durable agent.

> **Verify you're on the `tailnet`**
>
> Run the following command:
> ```bash
> tailscale status
> ```
>
> If you see **Logged Out** then you need to reauthenticate to the `tailnet`
>
> Run the following command to authenticate to the `tailnet`
> ```bash
> tailscale up --auth-key="$TS_AUTHKEY" --hostname="${WORKSHOP_USER_ID}-env"
> ```

## Environment

All code for this exercise lives in `exercises/04_go_agent/practice/`. Unlike the earlier exercises, this one has no **TODO**s, the code is already complete. You'll read, run, and optionally tweak it in place. Step 5 invites you to customize the Claude prompt in `activities.go` if you want to experiment.

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
└──────────────────┘            │ shared LLM gateway)│   │ (shared key)    │
                                └────────────────────┘   └─────────────────┘
```

The Worker joins the `tailnet` *itself* via `tsnet` (just like Exercise 2), but this time its HTTP client reaches a real `metrics-server` node and its LLM calls go through Aperture to Claude. All three destinations are on the workshop `tailnet`, with no public endpoints.

## Step 1: Scavenger hunt through the code

No TODOs this exercise, so before you run anything take a few minutes in the **Code Editor** tab familiarizing yourself with the code and see if you can find the answer to these questions . Treat it as orienteering, not a quiz.

- What is the default value of `HEALTH_CHECK_INTERVAL` if you don't set one?
- How many times does the Schedule fire before it pauses itself?
- Which two activities does `HealthCheckWorkflow` run, and in what order?
- What data flows from the first activity into the second?
- Where does the HTTP traffic actually go before it reaches Anthropic?
- What is the default `AI_MODEL`?

Keep your answers in your head. You'll see most of these values in the Worker's log output and the Temporal UI over the next few steps.

## Step 2: Start the Worker

The Worker is one process. The starter (which creates the Schedule) is a separate process. That split mirrors how you'd run this in production: many Workers polling a task queue, a separate tool registering and managing Schedules against the Temporal Server.

Start the Worker first. It joins the `tailnet` via `tsnet`, dials Temporal through that node, and begins polling its task queue for Workflow tasks. It will sit idle until the starter creates a Schedule that fires Workflows onto that queue.

In the **Worker** terminal:

```bash
cd exercises/04_go_agent/practice
export METRICS_URL=http://metrics-server:9100/metrics
go run . worker
```

`WORKSHOP_USER_ID`, `TS_AUTHKEY`, and `APERTURE_URL` are already exported by the workshop setup, so `METRICS_URL` is the only variable you need to set.

The first run takes 10 to 30 seconds while `tsnet` registers the node. You should see log lines like:

```output
level=INFO msg="joined tailnet" hostname=<you>-ex4-metrics-worker-<5 random chars> userID=<you>
level=INFO msg="connected to temporal" host=temporal-dev:7233
level=INFO msg="metrics reachable" url=http://metrics-server:9100/metrics
level=INFO msg="worker running" taskQueue=<you>-health-check
```

The Worker is now polling and waiting for work.

## Step 3: Create the Schedule

The starter is a short-lived process that deletes any existing Schedule with your ID and then creates a new one with `TriggerImmediately: true`. That immediate trigger is what fires the first Workflow run; subsequent runs fire on the interval.

In the **Starter** terminal:

```bash
cd exercises/04_go_agent/practice
export HEALTH_CHECK_INTERVAL=1m
go run . starter
```

`HEALTH_CHECK_INTERVAL` controls how often the Schedule fires; you'll change it in Step 5 to see the Schedule react. The starter joins the `tailnet` as its own node, connects to Temporal, creates the Schedule, and exits. You should see:

```output
level=INFO msg="joined tailnet" hostname=<you>-ex4-metrics-starter-<5 random chars> userID=<you>
level=INFO msg="connected to temporal" host=temporal-dev:7233
level=INFO msg="created schedule" id=<you>-health-check-schedule interval=1m0s workflow=<you>-health-check remainingActions=5
```

Back in the **Worker** terminal you should now see the Worker pick up the first fired Workflow and run the two activities (`FetchMetrics` then `AnalyzeMetrics`).

## Step 4: Watch the Schedule in the Temporal UI

Open the **Temporal UI** tab. Two places to look:

- **Schedules**, then `<your-user-id>-health-check-schedule`. Shows the interval, the next fire time, and recent runs.
- **Workflows**, then search for `<your-user-id>-health-check`. Each fired run has a completed row (its ID is suffixed with the Schedule fire time) whose Result panel contains the `HealthReport` JSON.

You should see at least one completed Workflow run whose result is a structured `HealthReport` produced by Claude.

> **Note:** If the **Temporal UI** tab shows a connection error or stale content, click the refresh button at the top of the tab. The iframe can hold an old render from before the `tailnet` was ready.

> **Note:** The Schedule is capped at 5 runs to keep the shared Temporal server clean. Once all 5 fire, the Schedule pauses itself. Schedules live on the Temporal Server, not on the Worker, so restarting the Worker does nothing to the Schedule. To reset the count, re-run the starter, which deletes and recreates the Schedule.

## Step 5: Tune the cadence

To change the interval, re-run the starter with a different `HEALTH_CHECK_INTERVAL`. The Worker keeps running; only the Schedule changes.

In the **Starter** terminal:

```bash
export HEALTH_CHECK_INTERVAL=30s
go run . starter
```

Any Go duration works (`30s`, `2m`, `5m`). The starter deletes the old Schedule and creates a new one with the new interval, so the 5-run count resets too. The Worker in the other terminal notices the new fires immediately.

## Step 6: Customize the Claude prompt (optional)

If you want to see how the summary changes when you change what you ask Claude, you can edit the prompt directly and restart the Worker.

Open `activities.go` in the **Code Editor** tab, find `AnalyzeMetrics`. The prompt lives in a raw string. Change it however you like, ask Claude to flag anything unusual, add a field to the `HealthReport` struct, or try a different tone. Then restart the Worker (the starter doesn't need restarting; the Schedule is unchanged) and watch the next Schedule fire produce a different `HealthReport` in the UI.

## Step 7: Run the offline tests

The Workflow and activities also come with offline tests that mock `node_exporter` and Aperture with `httptest.Server`, so they don't need the `tailnet` at all.

In the **Worker** terminal (after stopping the Worker with `Ctrl+C`):

```bash
go test ./...
```

You should see all tests pass. These are the tests that would run in CI for a production version of this service.

## Wrapping Up

In this exercise you:

- Ran a Go Worker that joins the `tailnet` via `tsnet` and dials three different `tailnet` services (Temporal, `node_exporter`, Aperture) through the same embedded node
- Used the same Aperture pattern from Exercise 3, this time with Anthropic's Claude instead of OpenAI
- Ran a separate starter process that registered a Temporal Schedule with `TriggerImmediately`, matching the production pattern of decoupling Workers from Schedule management
- Watched the Schedule fire on creation and on a cadence in the Temporal UI, with the Worker picking up each fired run
- Tuned the cadence by re-running the starter with a different `HEALTH_CHECK_INTERVAL`, without restarting the Worker
- Optionally customized the Claude prompt and saw the structured `HealthReport` change on the next fire
- Ran the offline tests that mock `node_exporter` and Aperture, no `tailnet` required

## Workshop Conclusion

Over four exercises you assembled a durable AI agent piece by piece and secured every connection it makes on a Tailscale `tailnet`, with zero public-internet exposure:

1. **Exercise 1** connected a Python Worker in a GCP VM to a Temporal dev server on a DigitalOcean VPS using the `tailnet` as the only network path, and declared that connection with a `temporal.toml` profile so the same code works in any environment.
2. **Exercise 2** moved the Worker itself onto the `tailnet` via `tsnet`. The Worker became its own node instead of riding on a system Tailscale client. You proved this by taking the system `tailscale` binary offline and watching the Go Worker join the `tailnet` on its own.
3. **Exercise 3** put an LLM inside a Temporal Workflow. OpenAI calls routed through **Aperture**, a `tailnet`-only API gateway that uses your Tailscale identity for auth and rate limiting instead of client-side API keys. You turned a single tool-call Workflow into a full agentic loop where the LLM chose its own tools at runtime and Temporal persisted every decision in the Workflow history.
4. **Exercise 4** combined both patterns in Go: a `tsnet`-embedded Worker that runs an agent on a Temporal Schedule, calling Claude through Aperture and scraping metrics off a `tailnet`-only service.

Three patterns to take home:

- **Temporal Client profiles** decouple connection settings from code. Swap environments with an env var instead of a code change.
- **`tsnet`-embedded Workers** let you deploy Workers anywhere with no network wiring. Each Worker is a first-class `tailnet` node.
- **Durable AI agents on Temporal** are real. The LLM picks tools at runtime, Temporal runs them as dynamic activities, and the Workflow history captures every decision so a crash mid-reasoning resumes exactly where it left off.

No ingress, no egress, no API keys on your machine. Every connection in every exercise went over the `tailnet`.
