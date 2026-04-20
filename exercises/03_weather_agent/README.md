# Exercise 3: Python Weather Agent

Build a durable AI agent that chains tool calls — with all LLM requests secured through Aperture.

## Goal

Run a Temporal workflow that uses OpenAI's function calling to autonomously answer questions. The LLM decides which tools to call, Temporal executes them as activities, and Aperture proxies every LLM request through the Tailscale network.

This exercise has two phases:
- **Phase A**: A simple tool-calling workflow (single LLM decision)
- **Phase B**: A full agentic loop (the LLM reasons through multiple steps)

## Background

### The Tool-Calling Pattern

The `ToolCallingWorkflow` makes one LLM call with a weather alerts tool available. If the LLM decides to use it, the workflow executes the tool as an activity and makes a second LLM call to format the result.

### The Agentic Loop Pattern

The `AgentWorkflow` is more powerful. It runs in a loop:

1. Ask the LLM what to do (with all tools available)
2. If the LLM picks a tool → execute it as a dynamic Temporal activity
3. Feed the result back to the LLM
4. Repeat until the LLM has enough information to respond

Example: "What's the weather where I am?" triggers:
- `get_ip_address` → gets your public IP
- `get_location_info` → geolocates the IP to a city/state
- `get_weather_alerts` → gets NWS alerts for that state
- Final response with the weather information

All of this is durable — if the worker crashes mid-loop, Temporal replays and continues.

### How Aperture Fits In

Every LLM call goes through Aperture instead of directly to OpenAI. Aperture:
- Holds the shared OpenAI API key (you don't need one)
- Identifies you by your Tailscale identity
- Enforces per-user rate limits

You'll configure this by setting `OPENAI_BASE_URL` to the Aperture endpoint.

### Connection

The worker and starter already use `ClientConfig.load_client_connect_config()` — the same environment configuration you set up in Exercise 1. No connection changes needed.

## Phase A: Tool-Calling (TODO 1)

### TODO 1 — Route LLM calls through Aperture

Open `practice/activities.py`. Find TODO 1 (~line 34) and add the Aperture base URL:

```python
client = AsyncOpenAI(
    max_retries=0,
    base_url=os.getenv("OPENAI_BASE_URL"),
)
```

This tells the OpenAI client to send requests to Aperture instead of directly to `api.openai.com`. Aperture proxies the request to OpenAI with the shared API key and tracks your usage via your Tailscale identity.

### Run Phase A

Terminal 1 — start the worker:
```bash
cd exercises/03_weather_agent/practice
uv run worker.py
```

Terminal 2 — run the workflow:
```bash
cd exercises/03_weather_agent/practice
uv run starter.py "What are the weather alerts in California?"
```

You should see the LLM call the weather alerts tool and return formatted results. Check the Temporal UI at `http://temporal-dev:8233` to see the workflow execution.

## Phase B: Agentic Loop (TODOs 2-3)

### TODO 2 — Enable the loop

Open `practice/agent_workflow.py`. Find TODO 2 (~line 22) and change `False` to `True`:

```python
while True:
```

This lets the agent make multiple tool calls in sequence until it has enough information to answer.

### TODO 3 — Execute the dynamic activity

In the same file, find TODO 3 (~line 65) and replace the empty string with the activity execution:

```python
tool_result = await workflow.execute_activity(
    item.name,
    args,
    start_to_close_timeout=timedelta(seconds=30),
)
```

This is the key line — `item.name` is the tool name the LLM chose (like `"get_ip_address"`), and Temporal executes it as a dynamic activity. The LLM decided what to call, Temporal makes it durable.

### Run Phase B

Terminal 1 — start the agent worker:
```bash
cd exercises/03_weather_agent/practice
uv run worker.py --agent
```

Terminal 2 — run the agent:
```bash
cd exercises/03_weather_agent/practice
uv run starter.py --agent "What's the weather like where I am right now?"
```

Watch the worker logs — you'll see the LLM chain through multiple tools before responding. Check the Temporal UI to see each tool call as a separate activity in the workflow history.

## Stuck?

Check the `solution/` directory for the completed code.
