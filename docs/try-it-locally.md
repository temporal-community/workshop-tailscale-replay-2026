# Try it locally

For a single developer who wants to run the workshop on their own machine, after attending or just out of curiosity. You'll run the full Tailscale + Temporal stack. A `temporal-ts-net` process joins your personal tailnet as `temporal-dev`, and a worker on the same machine talks to it over the mesh. That's the whole point of the workshop; we keep it intact.

The only piece we swap is Aperture. Aperture is Tailscale-internal, so a solo setup uses a direct OpenAI key instead. Everything else works the same.

If you're planning to **teach** the workshop (whether via Instruqt or without), see [Running with Instruqt](run-with-instruqt.md) or [Running without Instruqt](run-without-instruqt.md) instead.

## Prerequisites

- Python 3.13 and `uv`
- Go 1.26+ (needed to build `temporal-ts-net`, and for Exercises 2 and 4)
- The Temporal CLI (`brew install temporal`)
- The [Tailscale client](https://tailscale.com/download) and a free Tailscale account
- An OpenAI API key

## Set up your tailnet

If you don't already have one, sign in at [login.tailscale.com](https://login.tailscale.com) with a GitHub/Google account. The free tier is plenty.

Install the Tailscale client on your machine and bring it up once so your Mac is on the tailnet:

```shell
tailscale up
```

Then generate **two auth keys** in the [admin console](https://login.tailscale.com/admin/settings/keys), both **reusable** and **pre-authorized**:

1. One for `temporal-ts-net` (the Temporal server).
2. One for re-authing your machine if needed (optional; your system Tailscale is already up).

Write them down somewhere; you'll use them below.

## Install the workshop repo

```shell
git clone https://github.com/temporal-community/workshop-tailscale-replay-2026
cd workshop-tailscale-replay-2026
uv sync
```

## Install `temporal-ts-net`

This is the Temporal CLI extension that puts the dev server on your tailnet. Builds from source:

```shell
git clone https://github.com/temporal-community/temporal-ts-net
cd temporal-ts-net
go install ./cmd/temporal-ts_net
cd ..
temporal help --all | grep ts-net   # confirm the extension is picked up
```

## Configure your Temporal profile

Create `~/.config/temporalio/temporal.toml`. Note the server address is `temporal-dev:7233`, reached via the tailnet, not `localhost`:

```toml
[profile.local]
address = "temporal-dev:7233"
namespace = "default"

[profile.local.env]
OPENAI_API_KEY = "sk-..."
OPENAI_BASE_URL = "https://api.openai.com/v1"
```

In every terminal you open: `export TEMPORAL_PROFILE=local`.

## Run the stack

Three terminals, all on your machine.

**Terminal 1, `temporal-ts-net`.** Joins your tailnet as `temporal-dev` and runs the dev server behind it:

```shell
TS_AUTHKEY=tskey-auth-<key-1> temporal ts-net \
    --db-filename ~/.temporal/workshop.db
```

First run takes a few seconds while it authenticates. When you see `Tailnet gRPC: temporal-dev:7233`, it's ready.

**Terminal 2, worker.** Your Mac's system Tailscale client gives this terminal tailnet connectivity; the worker talks to `temporal-dev:7233` as if it were on a local network:

```shell
export TEMPORAL_PROFILE=local
cd exercises/03_weather_agent/practice
uv run worker.py
```

**Terminal 3, starter.**

```shell
export TEMPORAL_PROFILE=local
cd exercises/03_weather_agent/practice
uv run starter.py "What's the weather where I am?"
```

Open the UI at `http://temporal-dev:8233` (again, the tailnet hostname, not `localhost`). You should see your workflow, and its activities resolving through the agentic loop.

## What's different vs. the real workshop

| Piece | Real workshop | Solo |
|---|---|---|
| Temporal server | Shared VPS, one per cohort | Your machine, runs `temporal-ts-net` locally |
| Tailnet | Workshop tailnet | **Your personal tailnet** |
| Aperture | Holds the shared OpenAI key, per-identity rate limits | **Swapped for direct OpenAI** - you provide your own key |
| Web UI | `http://temporal-dev:8233`, the whole room sees it | `http://temporal-dev:8233`, only you see it |

You skip the rate-limit demo (that's Aperture's party trick) but get everything else, including the `temporal-ts-net` setup, the tailnet addressing, and the full durable-agent story.

!!! tip "Why not just use `temporal server start-dev` on localhost?"
    Because then you're not running the workshop - you're running a Temporal agent. The point of this workshop is `temporal-ts-net` and Tailscale. Swapping `temporal-dev:7233` for `localhost:7233` removes the networking story that the exercises build toward.
