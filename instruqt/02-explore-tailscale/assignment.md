---
slug: explore-tailscale
type: challenge
title: "Explore Your Tailscale Network"
teaser: Discover what's on the tailnet and understand how Aperture secures API calls
difficulty: basic
timelimit: 900
notes:
  - type: text
    contents: |-
      # Exploring the Network

      Now that you've proven the tailnet works, let's look under the hood.
      You'll discover all the machines on the network and learn how
      Aperture acts as a security gateway for LLM API calls.
tabs:
  - title: Terminal
    type: terminal
    hostname: workshop
  - title: Temporal UI
    type: service
    hostname: workshop
    port: 8233
    url: http://temporal-dev:8233
---

# Exercise 2: Exploring Your Tailscale Network

Now that you've run a workflow through the tailnet, let's understand the network you're on.

## Part 1: Discover Your Network

### See all machines on the tailnet

```bash
tailscale status
```

You should see:
- **Your machine** — the Instruqt VM you're working on
- **temporal-dev** — the VPS running the shared Temporal dev server
- **Other attendee machines** — everyone else in the workshop
- **Aperture endpoint** — the API gateway that will proxy your LLM calls

### Ping the Temporal server

```bash
tailscale ping temporal-dev
```

Notice the latency. This is a direct, encrypted WireGuard connection.

### Check your Tailscale identity

```bash
tailscale whois $(tailscale ip -4)
```

This shows your identity on the tailnet. Aperture uses this identity to track your API usage — no API keys needed on your machine.

### Verify what's accessible

```bash
# gRPC port
nc -zv temporal-dev 7233

# Web UI
curl -s -o /dev/null -w "%{http_code}" http://temporal-dev:8233
```

## Part 2: Understanding Aperture

<!-- ============================================================
     KARTIK: This section is yours to fill.

     Suggested content:
     - How Aperture proxies API calls (architecture)
     - How it identifies users (by Tailscale identity)
     - Rate limiting configuration for this workshop
     - Quick demo of the Aperture dashboard/metrics
     - Any Tailscale Serve / Funnel demonstration

     Target time: ~10 minutes
     ============================================================ -->

### How Aperture secures your LLM calls

*[Kartik to fill: explanation of Aperture's role as an API gateway]*

### Try the Aperture endpoint

```bash
# Example — update with actual endpoint from Kartik:
# curl https://aperture.<tailnet-name>.ts.net/v1/models
```

## What You've Learned

- How to discover machines on a Tailscale network
- That Tailscale provides identity for free — no extra auth needed
- How Aperture uses that identity to proxy and rate-limit API calls
- The full network path: Your VM → Tailscale → Temporal Server / Aperture → OpenAI

Next up: you'll use Aperture to power an AI agent workflow.
