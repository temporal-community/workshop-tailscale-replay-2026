# Workshop overview

A 90-minute, hands-on workshop. Attendees end the session with a durable AI agent running on a shared Temporal server, routing its LLM calls through Aperture, all over an encrypted Tailscale mesh.

## Learning objectives

By the end of the workshop, an attendee can:

1. Explain why durable execution matters for AI agents (retries, crashes, multi-step reasoning).
2. Connect a worker to a remote Temporal server over Tailscale without any VPN setup.
3. Route LLM calls through an API gateway and understand how identity-aware rate limiting protects a shared key.
4. Enable an agentic tool-calling loop in a Temporal workflow and watch the decisions unfold in the UI.

## What attendees build

A weather agent. It answers questions like *"what's the weather where I am?"* by chaining three tool calls:

1. `get_ip_address` - finds the caller's public IP.
2. `get_location_info` - geocodes that IP.
3. `get_weather_alerts` - queries the National Weather Service.

The LLM picks each tool, Temporal executes it as an Activity, the result is fed back to the LLM, and the loop continues until the model decides it has enough information to answer.

Every LLM call goes through Aperture, which holds the real OpenAI key and applies per-user quotas. Everything talks over a Tailscale mesh. Attendees never open a firewall or hand-configure a VPN.

## Time budget

| Segment | Length |
|---|---|
| Intro + architecture walk-through | 10 min |
| Exercise 1 - Hello Tailnet | 15 min |
| Exercise 2 - Explore Tailscale + Aperture | 15 min |
| AI agents on Temporal (talk) | 10 min |
| Exercise 3 - Weather agent | 25 min |
| Rate-limit demo (everyone fires at once) | 5 min |
| `temporal-ts-net` + Go agent preview | 5 min |
| Wrap-up + Q&A | 5 min |

Exercise 4 (the Go agent) is a take-home stretch goal. Files are stubbed and ready.

## What attendees need

- An Instruqt VM (or the [without-Instruqt](run-without-instruqt.md) equivalent) with Python 3.13, `uv`, Go 1.26+, and the Tailscale client pre-installed.
- A browser for the Temporal Web UI at `http://temporal-dev:8233`.
- **No OpenAI key.** Aperture handles that.
- **No Tailscale account.** The workshop tailnet is pre-provisioned; the VM joins via a reusable auth key.

## Prerequisites for attendees

- Comfortable reading Python. They modify existing code, not write from scratch. Every exercise is a TODO inside a complete scaffold.
- Light familiarity with Temporal helps, but isn't required. Exercise 1 is intentionally a gimme to build muscle memory.
- No prior Tailscale or Aperture experience assumed.
