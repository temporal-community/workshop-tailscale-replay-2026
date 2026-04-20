# Exercise 2: Exploring Your Tailscale Network

Now that you've run a workflow through the tailnet, let's understand the network you're on.

## Goal

Use Tailscale CLI tools to discover what's on the network, understand how access control works, and see how Aperture fits into the architecture.

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

Notice the latency. This is a direct, encrypted WireGuard connection — no relay servers, no VPN concentrators.

### Check what's accessible

Try reaching the Temporal gRPC port and Web UI:

```bash
# gRPC port (should succeed)
nc -zv temporal-dev 7233

# Web UI (should succeed — you already used this in Exercise 1)
curl -s -o /dev/null -w "%{http_code}" http://temporal-dev:8233
```

### Check your Tailscale identity

```bash
tailscale whois $(tailscale ip -4)
```

This shows your identity on the tailnet. Aperture uses this identity to track your API usage — no API keys needed on your machine.

## Part 2: Understanding Aperture

<!-- ============================================================
     KARTIK: This section is yours to fill.

     The attendees have:
     - Verified Tailscale connectivity (Exercise 1)
     - Explored the tailnet with CLI tools (Part 1 above)
     - Seen their Tailscale identity

     Suggested content for this section:
     - How Aperture proxies API calls (architecture)
     - How it identifies users (by Tailscale identity)
     - Rate limiting configuration for this workshop
     - Quick demo of the Aperture dashboard/metrics
     - Any Tailscale Serve / Funnel demonstration

     Target time: ~10 minutes of this 15-minute exercise.
     ============================================================ -->

### How Aperture secures your LLM calls

*[Kartik to fill: explanation of Aperture's role as an API gateway]*

### Aperture endpoint

The Aperture proxy is available at:

```
# TODO: Kartik to provide the actual endpoint URL
https://aperture.<tailnet-name>.ts.net/v1
```

This endpoint accepts OpenAI-compatible requests and forwards them to OpenAI with the shared API key. Your Tailscale identity is used to enforce per-user rate limits.

### Try it out

*[Kartik to fill: optional curl command to test Aperture directly]*

```bash
# Example (update with actual endpoint):
# curl https://aperture.<tailnet-name>.ts.net/v1/models \
#   -H "Authorization: Bearer <aperture-token>"
```

## What You've Learned

- How to discover machines on a Tailscale network
- That Tailscale provides identity for free — no extra auth layer needed
- How Aperture uses that identity to proxy and rate-limit API calls
- The full network path: Your VM → Tailscale → Temporal Server / Aperture → OpenAI

Next up: you'll use this Aperture endpoint to power an AI agent workflow.
