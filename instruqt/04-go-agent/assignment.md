---
slug: go-agent
type: challenge
title: "Go Agent (Stretch Goal)"
teaser: Port the weather agent to Go — same Temporal server, same Aperture, different language
difficulty: advanced
timelimit: 1800
notes:
  - type: text
    contents: |-
      # Stretch Goal: Go Agent

      This is a take-home challenge. The Go files are stubbed out with
      TODOs. Use the Python implementation as your reference and build
      the same weather agent in Go.
tabs:
  - title: Terminal
    type: terminal
    hostname: workshop
  - title: Code Editor
    type: code
    hostname: workshop
    path: /root/workshop/exercises/04_go_agent/practice
  - title: Temporal UI
    type: service
    hostname: workshop
    port: 8233
    url: http://temporal-dev:8233
---

# Exercise 4: Go Agent (Stretch Goal)

The same weather agent pattern, implemented in Go.

## Goal

Port the Python agentic loop to Go, demonstrating that Temporal + Tailscale + Aperture work across languages. Same shared Temporal server, same Aperture endpoint, different language.

## Architecture

The Go agent follows the same pattern:

1. `CreateCompletion` activity calls the OpenAI API through Aperture
2. The workflow loops: ask the LLM → execute chosen tool → feed result back
3. Tool activities (`GetWeatherAlerts`, `GetIPAddress`, `GetLocationInfo`) call the same public APIs
4. Worker connects to the shared Temporal server via the same `temporal.toml` config

The key difference: Go doesn't have an official OpenAI SDK integrated with Temporal, so you'll make HTTP requests directly to the Aperture endpoint.

## Files

Open the **Code Editor** tab to see the stubbed files:

| File | What to implement |
|------|-------------------|
| `main.go` | Worker setup and workflow starter |
| `workflow.go` | Agentic loop workflow |
| `activities.go` | OpenAI API call + tool implementations |

## Getting Started

```bash
cd /root/workshop/exercises/04_go_agent/practice
go mod tidy
go run .         # Start the worker
go run . run "What's the weather like where I am?"  # Run a workflow
```

## Hints

- The Aperture endpoint is `$OPENAI_BASE_URL/responses` (POST)
- Set `Authorization: Bearer $OPENAI_API_KEY` header
- Use `encoding/json` for request/response bodies
- Look at `exercises/03_weather_agent/solution/` for the Python reference

## Solution

The solution will be added in a future update. This is a take-home challenge!
