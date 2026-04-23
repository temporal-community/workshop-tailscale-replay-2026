---
slug: weather-agent
id: ubqh4feyxndc
type: challenge
title: Python Weather Agent
teaser: Build a durable AI agent with LLM calls routed through Aperture
notes:
- type: text
  contents: |-
    # AI Agent Time

    Now you'll build a durable AI weather agent. The LLM autonomously
    chains tool calls — get your IP, geolocate it, fetch weather alerts —
    all orchestrated by Temporal, all LLM calls secured through Aperture.
tabs:
- id: cd5f6dijn3vr
  title: Code Editor
  type: code
  hostname: workshop
  path: /root/workshop/exercises/03_weather_agent/practice
- id: dz8wzkrpjtxa
  title: Worker
  type: terminal
  hostname: workshop
  workdir: /root/workshop
- id: xnomxmss7f74
  title: Starter
  type: terminal
  hostname: workshop
  workdir: /root/workshop
- id: ryto8syprq37
  title: Temporal UI
  type: service
  hostname: workshop
  port: 8233
- id: rxzw3xyphu9f
  title: Aperture UI
  type: service
  hostname: workshop
  port: 80
difficulty: basic
timelimit: 1500
enhanced_loading: null
---

# Exercise 3: Python Weather Agent

Build a durable AI agent that chains tool calls — with all LLM requests secured through Aperture.

## Background

This exercise has two phases:

**Phase A: Tool-Calling** — A simple workflow where the LLM decides whether to call a weather tool.

**Phase B: Agentic Loop** — A full loop where the LLM reasons through multiple steps:
1. "What's the weather where I am?" → calls `get_ip_address`
2. Gets your IP → calls `get_location_info`
3. Gets your city/state → calls `get_weather_alerts`
4. Has enough info → responds with the weather

Every LLM call goes through **Aperture** instead of directly to OpenAI. Aperture holds the shared API key, identifies you by your Tailscale identity, and enforces rate limits.

> **Not on the tailnet?** If you joined late or `tailscale status` shows **Logged out**, run this in the **Worker** terminal first:
>
> ```bash
> tailscale up --auth-key="$TS_AUTHKEY" --hostname="${WORKSHOP_USER_ID}-env"
> ```

## Phase A: Tool-Calling

### **TODO 1** - Route LLM calls through Aperture

Open `activities.py` in the **Code Editor** tab. Find **TODO 1** and add the Aperture base URL to the OpenAI client:

```python
client = AsyncOpenAI(
    max_retries=0,
    base_url=f"{os.getenv('APERTURE_URL')}/v1",
    api_key="",  # Aperture ignores this; identity comes from Tailscale WhoIs
)
```

This tells the OpenAI client to send requests to Aperture instead of `api.openai.com`.

### Run Phase A

In the **Terminal** tab, start the worker:

```bash
cd /root/workshop/exercises/03_weather_agent/practice
uv run worker.py
```

Open a **new terminal** (`+` button) and run the workflow:

```bash
cd /root/workshop/exercises/03_weather_agent/practice
uv run starter.py "What are the weather alerts in California?"
```

You should see the LLM call the weather tool and return results. Check the **Temporal UI** tab to see the workflow.

## Phase B: Agentic Loop

### **TODO 2** - Enable the loop

Open `agent_workflow.py` in the Code Editor. Find **TODO 2** and change `False` to `True`:

```python
while True:
```

### **TODO 3** - Execute the dynamic activity

In the same file, find **TODO 3** and replace the empty string:

```python
tool_result = await workflow.execute_activity(
    item.name,
    args,
    start_to_close_timeout=timedelta(seconds=30),
)
```

`item.name` is the tool the LLM chose (like `"get_ip_address"`), and Temporal executes it as a dynamic activity.

### Run Phase B

Stop the previous worker (Ctrl+C), then start the agent worker:

```bash
cd /root/workshop/exercises/03_weather_agent/practice
uv run worker.py --agent
```

In the other terminal:

```bash
cd /root/workshop/exercises/03_weather_agent/practice
uv run starter.py --agent "What's the weather like where I am right now?"
```

Watch the worker logs — the LLM chains through multiple tools before responding. Check the **Temporal UI** to see each tool call as a separate activity.

## Stuck?

The solution files are in `/root/workshop/exercises/03_weather_agent/solution/`.
