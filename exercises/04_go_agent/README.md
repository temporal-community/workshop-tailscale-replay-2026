# Exercise 4: Go Agent (Stretch Goal)

The same weather agent pattern, implemented in Go.

## Goal

Port the Python agentic loop to Go, demonstrating that Temporal + Tailscale + Aperture work across languages. The same shared Temporal server, the same Aperture endpoint, different language.

## Status

This exercise is a **take-home stretch goal**. The Go files are stubbed out with TODOs describing what each function should do. Use the Python implementation in Exercise 3 as your reference.

## Architecture

The Go agent follows the same pattern as the Python version:

1. `CreateCompletion` activity calls the OpenAI API through Aperture
2. The workflow loops: ask the LLM → execute chosen tool → feed result back
3. Tool activities (`GetWeatherAlerts`, `GetIPAddress`, `GetLocationInfo`) call the same public APIs
4. Worker connects to the shared Temporal server on the tailnet

The key difference: Go doesn't have an official OpenAI SDK integrated with Temporal, so you'll make HTTP requests directly to the Aperture endpoint (which is OpenAI-compatible).

## Files

| File | What to implement |
|------|-------------------|
| `practice/main.go` | Worker setup and workflow starter |
| `practice/workflow.go` | Agentic loop workflow |
| `practice/activities.go` | OpenAI API call + tool implementations |

## Getting Started

```bash
cd exercises/04_go_agent/practice
go mod tidy
go run . # Start the worker
go run . run "What's the weather like where I am?" # Run a workflow
```

## Hints

- The OpenAI Responses API endpoint is `POST /v1/responses`
- Set the `Authorization` header to `Bearer $OPENAI_API_KEY`
- Use `encoding/json` to marshal/unmarshal request/response bodies
- The `base_url` should come from `OPENAI_BASE_URL` environment variable
- Look at `exercises/03_weather_agent/solution/` for the complete Python reference

## Solution

The solution will be added in a future update. For now, use the Python implementation as your guide.
