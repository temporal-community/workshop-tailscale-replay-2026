# Securing AI Applications with Tailscale and Temporal

Workshop materials for the joint Tailscale + Temporal session at **Replay 2026**.

## What You'll Build

A durable AI weather agent powered by Temporal, secured by Tailscale:

- **Temporal** orchestrates an agentic loop where an LLM autonomously chains tool calls
- **Tailscale** provides zero-config encrypted networking between your machine and the shared infrastructure
- **Aperture** acts as an API gateway, proxying LLM calls with rate limiting and shared key management

```
┌─────────────────┐     Tailscale      ┌──────────────────┐
│  Your Instruqt  │◄──── Tailnet ─────►│  Temporal Dev     │
│  VM (Worker)    │                    │  Server (VPS)     │
│                 │                    │  temporal-dev:7233│
│  - Python agent │     Tailscale      │  temporal-dev:8233│
│  - Go agent     │◄──── Tailnet ─────►│                  │
│    (stretch)    │                    ├──────────────────┤
└─────────────────┘                    │  Aperture        │
                                       │  (API Gateway)   │
                                       │       │          │
                                       └───────┼──────────┘
                                               │
                                               ▼
                                       ┌──────────────────┐
                                       │  OpenAI API      │
                                       │  (shared key)    │
                                       └──────────────────┘
```

## Prerequisites

Your Instruqt VM comes pre-configured with everything you need:
- Python 3.13 and `uv`
- Tailscale (connected to the workshop tailnet)
- Workshop code cloned and dependencies installed

## Quick Start

Verify your environment:

```bash
uv run scripts/verify_setup.py
```

Then start with [Exercise 1](exercises/01_hello_tailnet/README.md).

## Exercises

| # | Exercise | Time | Description |
|---|----------|------|-------------|
| 1 | [Hello Tailnet](exercises/01_hello_tailnet/README.md) | 15 min | Run a geo-IP workflow on the shared Temporal server via Tailscale |
| 2 | [Explore Tailscale](exercises/02_explore_tailscale/README.md) | 15 min | Discover your network, understand Aperture, explore access control |
| 3 | [Weather Agent](exercises/03_weather_agent/README.md) | 25 min | Build a durable AI agent with LLM calls routed through Aperture |
| 4 | [Go Agent](exercises/04_go_agent/README.md) | Stretch | Same agent pattern in Go (take-home challenge) |

## Temporal Web UI

Once connected to the tailnet, open the shared Temporal Web UI:

```
http://temporal-dev:8233
```

You'll see everyone's workflows running on the shared server.

## Resources

- [temporal-ts-net](https://github.com/temporal-community/temporal-ts-net) — Temporal CLI extension for Tailscale
- [Temporal Python SDK](https://docs.temporal.io/develop/python)
- [Temporal Go SDK](https://docs.temporal.io/develop/go)
- [Tailscale Documentation](https://tailscale.com/kb)
- [Aperture Documentation](https://docs.tailscale.com/aperture)

## Instructor Setup

See [instructor/README.md](instructor/README.md) for VPS and Instruqt configuration.
