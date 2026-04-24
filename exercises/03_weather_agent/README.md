# Exercise 3: Python Weather Agent

Build a durable AI agent that chains tool calls, with all LLM requests secured through Aperture.

## Goal

Run a Temporal Workflow that uses OpenAI's function calling to autonomously answer questions. The LLM decides which tools to call, Temporal executes them as activities, and Aperture proxies every LLM request through the Tailscale network.

This exercise has two phases:
- **Phase A**: A simple tool-calling Workflow (single LLM decision)
- **Phase B**: A full agentic loop (the LLM reasons through multiple steps)

## Background

### The Tool-Calling Pattern

The `ToolCallingWorkflow` makes one LLM call with a weather alerts tool available. If the LLM decides to use it, the Workflow executes the tool as an activity and makes a second LLM call to format the result.

### The Agentic Loop Pattern

The `AgentWorkflow` is more powerful. It runs in a loop:

1. Ask the LLM what to do (with all tools available)
2. If the LLM picks a tool, execute it as a dynamic Temporal activity
3. Feed the result back to the LLM
4. Repeat until the LLM has enough information to respond

Example: "What's the weather where I am?" triggers:
- `get_ip_address`, which gets your public IP
- `get_location_info`, which geolocates the IP to a city and state
- `get_weather_alerts`, which gets NWS alerts for that state
- Final response with the weather information

All of this is durable. If the Worker crashes mid-loop, Temporal replays and continues.

### How Aperture Fits In

Every LLM call goes through Aperture instead of directly to OpenAI. Aperture:
- Holds the shared OpenAI API key (you don't need one)
- Identifies you by your Tailscale identity
- Enforces per-user rate limits

You'll configure this by pointing the OpenAI client at the `APERTURE_URL` endpoint (and passing a throwaway key that Aperture strips).

### Connection

The worker and starter already call `ClientConfig.load_client_connect_config(profile="tailnet")`, which reads the `tailnet` profile from the `temporal.toml` you inspected in Exercise 1. No connection changes needed.

## Phase A: Tool-Calling (TODO 1)

### TODO 1: Route LLM calls through Aperture

Open `practice/activities.py`. Find TODO 1 (~line 34) and add the Aperture base URL:

```python
client = AsyncOpenAI(
    max_retries=0,
    base_url=f"{os.getenv('APERTURE_URL')}/v1",
    api_key="",  # Aperture ignores this; identity comes from Tailscale WhoIs
)
```

This tells the OpenAI client to send requests to Aperture instead of directly to `api.openai.com`. Aperture proxies the request to the upstream provider with the real shared key and tracks your usage via your Tailscale identity. The OpenAI SDK requires an `api_key` to be present, but the value is discarded by Aperture.

### Run Phase A

In the **Worker** terminal:
```bash
cd exercises/03_weather_agent/practice
uv run worker.py
```

In the **Starter** terminal:
```bash
cd exercises/03_weather_agent/practice
uv run starter.py "What are the weather alerts in California?"
```

You should see the LLM call the weather alerts tool and return formatted results. Check the Temporal UI to see the Workflow execution.

## Phase B: Agentic Loop (TODOs 2 and 3)

### TODO 2: Enable the loop

Open `practice/agent_workflow.py`. Find TODO 2 (~line 22) and change `False` to `True`:

```python
while True:
```

This lets the agent make multiple tool calls in sequence until it has enough information to answer.

### TODO 3: Execute the dynamic activity

In the same file, find TODO 3 (~line 65) and replace the empty string with the activity execution:

```python
tool_result = await workflow.execute_activity(
    item.name,
    args,
    start_to_close_timeout=timedelta(seconds=30),
)
```

This is the key line. `item.name` is the tool name the LLM chose (such as `"get_ip_address"`), and Temporal executes it as a dynamic activity. The LLM decides what to call, and Temporal makes it durable.

### Run Phase B

In the **Worker** terminal, stop the previous Worker with `Ctrl+C`, then:
```bash
cd exercises/03_weather_agent/practice
uv run worker.py --agent
```

In the **Starter** terminal:
```bash
cd exercises/03_weather_agent/practice
uv run starter.py --agent "What's the weather like where I am right now?"
```

Watch the Worker logs. You'll see the LLM chain through multiple tools before responding. Check the Temporal UI to see each tool call as a separate activity in the Workflow history.

## Stuck?

Check the `solution/` directory for the completed code.
