---
theme: ./theme-temporal
title: Securing AI Applications with Tailscale and Temporal
info: |
  ## Replay 2026 Workshop
  A hands-on workshop on securing AI agent infrastructure with
  Tailscale networking and Aperture API gateway, powered by Temporal.
author: Mason Egger & Kartik Venugopal
keywords: temporal,tailscale,aperture,ai,agents
colorSchema: dark
fonts:
  sans: 'Inter'
  mono: 'Noto Sans Mono'
  weights: '200,300,400,500,600'
  italic: false
layout: cover
---

# Securing AI Applications with Tailscale and Temporal

## Replay 2026 Workshop

<br>

**Mason Egger** · Temporal<br>
**Kartik Bharath** · Tailscale

<!--
Welcome everyone! Over the next 90 minutes you'll build and run a durable AI agent
that's secured end-to-end by Tailscale networking.
-->

---
layout: two-cols
---

# About Us

### Mason Egger
Developer Advocate, **Temporal**

Builder of workshops, CLI tooling, and
educational content for the Temporal ecosystem.

::right::

<br><br>

### Kartik Bharath
*[Title]*, **Tailscale**

*[Kartik bio]*

<!--
Quick intros — we'll keep it brief since we have a lot to build today.
-->

---

# What We're Building Today

<br>

A **durable AI weather agent** that:

- 🔄 Uses **Temporal** to orchestrate an agentic tool-calling loop
- 🔒 Runs on a **Tailscale** private network — no public internet exposure
- 🚪 Routes all LLM calls through **Aperture** for rate limiting and key management
- 🌍 Works across everyone's machines — all hitting one shared Temporal server

<br>

> You will **not** need an OpenAI API key. Aperture handles that.

---

# The Problem

AI applications in production need more than just "call the LLM":

<v-clicks>

- **Durability** — What happens when your agent crashes mid-reasoning?
- **Networking** — How do distributed workers reach your infrastructure securely?
- **API Security** — How do you share expensive API keys without exposing them?
- **Rate Limiting** — How do you prevent one user from burning your entire budget?

</v-clicks>

<br>

<v-click>

Today we solve all four.

</v-click>

---
layout: section
---

# Architecture

---

# Architecture Overview

```mermaid {scale: 0.8}
flowchart LR
  subgraph VM["Your Instruqt VM"]
    PY[Python worker]
    GO[Go worker]
  end
  subgraph VPS["VPS"]
    TS["Temporal Dev Server<br/>temporal-dev:7233<br/>temporal-dev:8233"]
    AP["Aperture<br/>API Gateway"]
  end
  OAI["OpenAI API<br/>(shared key)"]
  VM <-. Tailnet .-> TS
  VM <-. Tailnet .-> AP
  AP --> OAI
```

<v-clicks>

- **Tailscale** — encrypted mesh network, zero config
- **Aperture** — API gateway with identity-based rate limiting
- **temporal-ts-net** — Temporal dev server exposed on the tailnet

</v-clicks>

---

# What is Tailscale?

<v-clicks>

- **Mesh VPN** built on WireGuard — every device connects directly
- **Zero config** — no firewall rules, no port forwarding, no VPN concentrators
- **Identity-based** — every connection knows who's on the other end
- **Tailnet** — your private network of devices

</v-clicks>

<br>

<v-click>

Your Instruqt VM is already connected. The shared Temporal server is just `temporal-dev:7233` — as if it were on your local network.

</v-click>

---

# What is Aperture?

<v-clicks>

- **API gateway** that sits between your code and external APIs
- **Shared key management** — one OpenAI key, many users, nobody sees the key
- **Identity-aware** — uses your Tailscale identity, no extra auth tokens
- **Rate limiting** — per-user quotas so no one burns the whole budget

</v-clicks>

<br>

<v-click>

Your LLM calls go to Aperture's endpoint instead of `api.openai.com`. Aperture forwards them with the real key and tracks your usage.

</v-click>

---

# temporal-ts-net

How the Temporal dev server got on the tailnet:

```bash
temporal ts-net \
    --db-filename /var/lib/temporal/workshop.db \
    --max-connections 2000 \
    --connection-rate-limit 200
```

<v-clicks>

- **Temporal CLI extension** — runs `temporal server start-dev` and proxies it onto the tailnet
- **No public exposure** — the server is only reachable via Tailscale
- **Supports gRPC + Web UI** — `temporal-dev:7233` and `temporal-dev:8233`
- **Built with tsnet** — Go library for embedding Tailscale in applications

</v-clicks>

---
layout: section
---

# Exercise 1: Hello Tailnet

---

# Exercise 1: What You'll Do

<br>

### Goal
Prove the tailnet works — run a workflow on the shared Temporal server.

<br>

### Steps

1. **Create `temporal.toml`** — configure the `tailnet` profile
2. **Add your user ID** to the workflow ID
3. **Start a worker** → it connects to `temporal-dev:7233`
4. **Run the workflow** → get your IP and geolocation
5. **Open the Temporal UI** → see your workflow alongside everyone else's

---

# Temporal Environment Configuration

Instead of hardcoding addresses, use a **config file**:

```toml
# ~/.config/temporalio/temporal.toml

[profile.default]
address = "localhost:7233"
namespace = "default"

[profile.tailnet]
address = "temporal-dev:7233"
namespace = "default"
```

<br>

The SDK reads this automatically:

```python
config = ClientConfig.load_client_connect_config()
client = await Client.connect(**config)
```

<v-click>

Set `TEMPORAL_PROFILE=tailnet` and every worker, starter, and CLI command connects to the right server.

</v-click>

---

# Exercise 1: Commands

<br>

### Copy the config

```bash
mkdir -p ~/.config/temporalio
cp temporal.toml.example ~/.config/temporalio/temporal.toml
```

### Start the worker

```bash
cd exercises/01_hello_tailnet/practice
uv run worker.py
```

### Run the workflow

```bash
cd exercises/01_hello_tailnet/practice
uv run starter.py
```

### Check the Temporal UI

Open **http://temporal-dev:8233** and find your workflow.

---
layout: exercise
heading: Exercise 1
minutes: 15
---

Create `temporal.toml`, add your user ID, run the geo-IP workflow.

Open the Temporal UI and find your workflow.

---
layout: section
---

# Exercise 2: Explore Tailscale

---

# Exercise 2: Your Tailscale Network

<br>

### Discover what's on the tailnet

```bash
tailscale status        # See all machines
tailscale ping temporal-dev   # Direct encrypted connection
tailscale whois $(tailscale ip -4)   # Your identity
```

<v-clicks>

- Your VM, the Temporal server, Aperture, and every other attendee
- Direct WireGuard connections — no relay servers
- Your identity is automatic — Aperture uses it for rate limiting

</v-clicks>

---

# How Aperture Secures Your LLM Calls

<!-- KARTIK: Replace this slide with your Aperture content -->

```mermaid {scale: 0.8}
sequenceDiagram
  autonumber
  participant VM as Your VM
  participant AP as Aperture
  participant OAI as OpenAI
  VM->>AP: POST /v1/responses<br/>(no API key needed)
  Note over AP: ✓ Identity: your-vm<br/>✓ Rate: 3 / 10 requests
  AP->>OAI: POST /v1/responses<br/>Authorization: Bearer sk-real-openai-key
  OAI-->>AP: response
  AP-->>VM: response
```

---
layout: exercise
heading: Exercise 2
minutes: 15
---

Explore your Tailscale network.

Run `tailscale status`, ping the server, check your identity.

---
layout: section
---

# AI Agents on Temporal

---

# The Tool-Calling Pattern

A single LLM decision: should I use a tool?

```mermaid {scale: 0.7}
flowchart TD
  U["User: 'Weather alerts in California?'"]
  LLM["LLM via Aperture"]
  A["Temporal Activity<br/>get_weather_alerts('CA')<br/>→ calls NWS API"]
  R["'There are 3 active alerts<br/>in California...'"]
  U --> LLM
  LLM -->|decides to call tool| A
  A -->|weather data| LLM
  LLM --> R
```

<v-click>

One LLM call decides, one activity executes, one final LLM call formats. Simple.

</v-click>

---

# The Agentic Loop Pattern

The LLM reasons through **multiple steps** autonomously:

```mermaid {scale: 0.7}
flowchart TD
  U["User: 'What's the weather where I am?'"]
  subgraph Loop["Agentic Loop · repeats until LLM is done"]
    direction LR
    P[LLM picks a tool] --> E[Execute Activity]
    E --> F[Feed result back to LLM]
    F --> P
  end
  U --> Loop
  Loop --> R[Final response]
```

<v-click>

`get_ip_address` → `get_location_info` → `get_weather_alerts` → respond

</v-click>

---

# Why This Needs Temporal

<v-clicks>

- **Each tool call is an Activity** — retried automatically on failure
- **The loop is a Workflow** — survives worker crashes, resumes from last completed step
- **Dynamic Activities** — the LLM picks the tool name, Temporal executes it
- **Durable state** — the entire conversation history is preserved

</v-clicks>

<br>

<v-click>

```python
# The LLM chose "get_ip_address" — Temporal runs it
tool_result = await workflow.execute_activity(
    item.name,  # dynamic — chosen by the LLM
    args,
    start_to_close_timeout=timedelta(seconds=30),
)
```

</v-click>

---

# How Aperture Fits the Agent

Every `create` activity call goes through Aperture:

```python
@activity.defn
async def create(request: OpenAIResponsesRequest) -> Response:
    client = AsyncOpenAI(
        max_retries=0,
        base_url=os.getenv("OPENAI_BASE_URL"),  # ← Aperture endpoint
    )
    return await client.responses.create(
        model=request.model,
        instructions=request.instructions,
        input=request.input,
        tools=request.tools,
    )
```

<v-click>

The tool execution activities (weather, IP, location) call **free public APIs** directly — only the LLM calls need Aperture.

</v-click>

---
layout: section
---

# Exercise 3: Weather Agent

---

# Exercise 3: Three TODOs

<br>

### TODO 1 — Route LLM calls through Aperture

`activities.py` — add `base_url` to the OpenAI client:

```python
client = AsyncOpenAI(
    max_retries=0,
    base_url=os.getenv("OPENAI_BASE_URL"),
)
```

<br>

Then run the **tool-calling workflow**:

```bash
uv run worker.py                    # Terminal 1
uv run starter.py "Weather alerts in California?"  # Terminal 2
```

---

# Exercise 3: Enable the Agentic Loop

<br>

### TODO 2 — Turn on the loop

`agent_workflow.py` — change `False` to `True`:

```python
while True:  # was: while False
```

### TODO 3 — Execute the dynamic activity

Same file — wire up the tool execution:

```python
tool_result = await workflow.execute_activity(
    item.name,
    args,
    start_to_close_timeout=timedelta(seconds=30),
)
```

---

# Exercise 3: Run the Agent

<br>

```bash
# Terminal 1 — start the agent worker
uv run worker.py --agent

# Terminal 2 — ask a question
uv run starter.py --agent "What's the weather like where I am?"
```

<br>

### What to watch for

- **Worker logs** — see the LLM chain: `get_ip_address` → `get_location_info` → `get_weather_alerts`
- **Temporal UI** — each tool call appears as a separate activity in the workflow history
- **The response** — a natural language answer with your local weather

---
layout: exercise
heading: Exercise 3
minutes: 25
---

Complete the 3 TODOs. Run the tool-calling workflow first, then enable the agentic loop.

Watch the multi-step reasoning in the Temporal UI.

---
layout: section
---

# Rate Limit Demo

---

# Let's All Fire at Once

<br>

Everyone run this at the same time:

```bash
cd exercises/03_weather_agent/practice
uv run starter.py --agent "What's the weather like where I am?"
```

<br>

<v-clicks>

- Watch the Aperture dashboard — per-user rate limits in action
- Some requests get throttled → Temporal **retries** the activity automatically
- Nobody's workflow fails — durability meets rate limiting

</v-clicks>

<!-- KARTIK: Show the Aperture dashboard here -->

---
layout: section
---

# temporal-ts-net & Go Agent

---

# How the Dev Server Got on the Tailnet

```go
// temporal-ts-net creates a tsnet.Server and proxies TCP connections
tsSrv := &tsnet.Server{
    Hostname: "temporal-dev",
    AuthKey:  os.Getenv("TS_AUTHKEY"),
}
tsSrv.Start()

// Listens on the tailnet, proxies to localhost:7233
listener, _ := tsSrv.Listen("tcp", ":7233")
for {
    conn, _ := listener.Accept()
    go proxy(conn, "localhost:7233")
}
```

<v-clicks>

- **6 lines** to put any TCP service on a Tailscale network
- Built as a **Temporal CLI extension** — `temporal ts-net`
- Supports rate limiting, max connections, idle timeouts
- Open source: [github.com/temporal-community/temporal-ts-net](https://github.com/temporal-community/temporal-ts-net)

</v-clicks>

---

# Go Agent Preview

**Exercise 4** — same weather agent, in Go. Take-home stretch goal.

```go
func AgentWorkflow(ctx workflow.Context, input string) (string, error) {
    for {
        // Call LLM through Aperture
        result := workflow.ExecuteActivity(ctx, CreateCompletion, ...)

        if result.Type == "function_call" {
            // Execute the tool the LLM chose
            toolResult := workflow.ExecuteActivity(ctx, result.Name, ...)
            // Feed result back to LLM
        } else {
            return result.Text, nil
        }
    }
}
```

Same Temporal server. Same Aperture endpoint. Same tailnet. Different language.

---
layout: section
---

# Wrap-Up

---

# What We Built

<br>

| Layer | Technology | What It Does |
|-------|-----------|--------------|
| **Durability** | Temporal | Orchestrates the agent loop, retries failures, survives crashes |
| **Networking** | Tailscale | Zero-config encrypted mesh between all machines |
| **API Security** | Aperture | Shared key management, identity-based rate limiting |
| **AI Agent** | OpenAI + Python | Multi-step reasoning with autonomous tool selection |

<br>

<v-click>

No VPN setup. No API keys on your machine. No hardcoded addresses.

Just a config file and a tailnet.

</v-click>

---

# Three Patterns to Take Home

<v-clicks>

### 1. Environment Configuration
Use `temporal.toml` profiles — not hardcoded addresses. Switch between local, staging, and production with one env var.

### 2. Aperture as API Gateway
Put expensive API keys behind a gateway with identity-based rate limiting. Your developers never see the key.

### 3. Durable AI Agents
Wrap your agentic loops in Temporal workflows. Every tool call is an activity. Every failure is a retry. Every crash is a resume.

</v-clicks>

---

# Resources

<br>

| Resource | Link |
|----------|------|
| Workshop repo | [github.com/temporal-community/workshop-tailscale-replay-2026](https://github.com/temporal-community/workshop-tailscale-replay-2026) |
| temporal-ts-net | [github.com/temporal-community/temporal-ts-net](https://github.com/temporal-community/temporal-ts-net) |
| Temporal Python SDK | [docs.temporal.io/develop/python](https://docs.temporal.io/develop/python) |
| Temporal Go SDK | [docs.temporal.io/develop/go](https://docs.temporal.io/develop/go) |
| Temporal envconfig | [docs.temporal.io/develop/environment-configuration](https://docs.temporal.io/develop/environment-configuration) |
| Tailscale docs | [tailscale.com/kb](https://tailscale.com/kb) |
| Aperture docs | [docs.tailscale.com/aperture](https://docs.tailscale.com/aperture) |

---
layout: end
---

# Questions?

**Mason Egger** · mason.egger@temporal.io

**Kartik Bharath** · *[email]*

<br>

Exercise 4 (Go Agent) is a take-home stretch goal — the files are stubbed and ready.
