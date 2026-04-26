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
    chains tool calls (get your IP, geolocate it, fetch weather alerts),
    all orchestrated by Temporal, with every LLM call secured through
    Aperture on the `tailnet`.
tabs:
- id: cd5f6dijn3vr
  title: Code Editor
  type: code
  hostname: workshop
  path: /root/workshop
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

Build a durable AI agent that chains tool calls, with all LLM requests secured through Aperture on the `tailnet`.

## Background

This exercise has two phases.

**Phase A: Tool-Calling.** A simple Workflow where the LLM decides whether to call a single weather tool.

**Phase B: Agentic Loop.** A full loop where the LLM reasons through multiple steps before answering:

1. "What's the weather where I am?", calls `get_ip_address`
2. Gets your IP, calls `get_location_info`
3. Gets your city and state, calls `get_weather_alerts`
4. Has enough info, responds with the weather

Every LLM call goes through **Aperture** instead of directly to OpenAI. Aperture holds the shared API key, identifies you by your Tailscale identity, and enforces rate limits. The same pattern applies in Exercise 4 with Anthropic's Claude.

## Environment

All code for this exercise lives in `exercises/03_weather_agent/`. Inside that directory:

- **`practice/`** is where you do your work. Each file has one or more **TODO** comments pointing at the change you need to make.
- **`solution/`** contains the finished version of every file. If you get stuck or want to double-check your work, compare against the matching file in `solution/`. Don't run from `solution/`, run from `practice/`.

> **Verify you're on the `tailnet`**
>
> In the [button label="Worker" background="#444CE7"](tab-1) terminal:
> ```bash,run
> tailscale status
> ```
>
> If you see **Logged Out**, reauthenticate:
> ```bash,run
> tailscale up --auth-key="$TS_AUTHKEY" --hostname="${WORKSHOP_USER_ID}-env"
> ```

## Step 1: Route LLM calls through Aperture

This step begins **Phase A: Tool-Calling**. The change is small but load-bearing. Instead of pointing the OpenAI client at `api.openai.com`, you point it at Aperture, the `tailnet`-only gateway that attaches your Tailscale identity, enforces rate limits, and swaps in the real API key server-side.

Open `exercises/03_weather_agent/practice/activities.py` in the [button label="Code Editor" background="#444CE7"](tab-0) tab. Find **TODO 1** and add the Aperture base URL to the OpenAI client:

```python
client = AsyncOpenAI(
    max_retries=0,
    base_url=f"{os.getenv('APERTURE_URL')}/v1",
    api_key="",  # Aperture ignores this; identity comes from Tailscale WhoIs
)
```

> **Heads up on paste indentation.** Python is whitespace-sensitive. If you copy the snippet above and paste over the existing `AsyncOpenAI(...)` call, double-check that every line inside the parens ends up at the same 4-space indent as the surrounding function body. A stray tab or extra spaces will throw an `IndentationError` at worker start.

`APERTURE_URL` is already exported in your environment. The OpenAI client will now send requests to Aperture, which forwards them to OpenAI with the shared API key after attaching your Tailscale identity.

## Step 2: Start the Phase A Worker

Now start the Worker. This is the process that executes the Workflow and its tool activities.

In the [button label="Worker" background="#444CE7"](tab-1) terminal:

```bash,run
cd exercises/03_weather_agent/practice
uv run worker.py
```

You should see the Worker connect to Temporal and start listening on its task queue.

## Step 3: Run the Phase A Workflow

With the Worker running, trigger the Workflow from the [button label="Starter" background="#444CE7"](tab-2) terminal.

```bash,run
cd exercises/03_weather_agent/practice
uv run starter.py "What are the weather alerts in California?"
```

You should see the LLM call the weather tool and return results. Click the [button label="Temporal UI" background="#444CE7"](tab-3) tab to find your Workflow and see the tool call execute as an activity.

> **Note:** If the **Temporal UI** tab shows a connection error or stale content, click the refresh button at the top of the tab. The iframe can hold an old render from before the `tailnet` was ready.

**What happened**

Your Worker executed a Workflow where the LLM chose to call a tool instead of answering directly. The LLM call itself flowed through Aperture, which authenticated you by your Tailscale identity and forwarded the request to OpenAI with the shared API key. Your code never touched an OpenAI API key.

## Step 4: Enable the agentic loop

This step begins **Phase B: Agentic Loop**. The difference from Phase A is that the Workflow now keeps handing the LLM tool results until the LLM decides it has enough information to answer, instead of stopping after one tool call.

Open `exercises/03_weather_agent/practice/agent_workflow.py` in the [button label="Code Editor" background="#444CE7"](tab-0) tab. Find **TODO 2** and change `False` to `True`:

```python
while True:
```

This turns the single-shot tool call from Phase A into a loop that repeatedly calls the LLM, executes whichever tool the LLM picked, feeds the result back, and stops only when the LLM decides it's done.

## Step 5: Execute the chosen activity dynamically

Still in `agent_workflow.py`, find **TODO 3** and replace the empty string with `item.name`:

```python
tool_result = await workflow.execute_activity(
    item.name,
    args,
    start_to_close_timeout=timedelta(seconds=30),
)
```

`item.name` is whichever tool the LLM picked on this iteration, such as `get_ip_address`, `get_location_info`, or `get_weather_alerts`. Temporal runs it as a dynamic activity, so the Worker does not hard-code which tool to call; Temporal dispatches by name.

## Step 6: Restart the Worker as the agent Worker

The Phase A Worker registered the single-shot tool-calling Workflow. For Phase B you need the agent Workflow registered, which means restarting the Worker with the `--agent` flag.

In the [button label="Worker" background="#444CE7"](tab-1) terminal, stop the previous Worker with `Ctrl+C`. You should still be in the `practice/` directory from Step 2, so start the agent Worker directly:

```bash,run
uv run worker.py --agent
```

You should see the Worker reconnect to Temporal and start listening on its task queue, this time with the agent Workflow registered.

## Step 7: Run the agentic Workflow

Now ask the agent a question that requires multiple tool calls to answer.

In the [button label="Starter" background="#444CE7"](tab-2) terminal:

```bash,run
uv run starter.py --agent "What's the weather like where I am right now?"
```

Watch the Worker logs. The LLM chains through multiple tools before responding: `get_ip_address`, then `get_location_info`, then `get_weather_alerts`, then a final answer. Click the [button label="Temporal UI" background="#444CE7"](tab-3) tab to see each tool call as a separate activity inside one Workflow execution.

**What happened**

The LLM made autonomous decisions about which tool to call next, and Temporal recorded every call, input, and output in the Workflow history. If the process had crashed halfway through, Temporal could replay the history on a new Worker and the agent would resume from exactly where it left off, even partway through a multi-tool reasoning chain.

## Step 8: Explore the Aperture UI

Open the [button label="Aperture UI" background="#444CE7"](tab-4) tab to see every LLM call your Workers made.

**Dashboard** gives you the aggregate view: total requests, total tokens, estimated cost, quota remaining, and a per-model breakdown in the **Metrics by Model** table. The **Recent Requests** list at the bottom shows individual calls with per-request token counts and costs.

Click the **Logs** tab to browse individual requests. Click any row to expand it and read the full request payload and response body for that call.

Click the **Tool Calls** tab to see every tool invocation that came out of the agentic loop, listed separately from the LLM calls that spawned them.

Click the **Adoption** tab for a cost and token-usage breakdown across models and over time.

> **Note:** Every Instruqt machine authenticated using the same `tag:infra`, so the Dashboard and Logs show requests from all attendees, not just yours. In a real deployment, Aperture attributes usage per user via their Tailscale identity from your IDP. Agentic workloads should have their own tags too — both so Aperture tracks them separately from human users and because zero-trust ACLs depend on a well-defined tag taxonomy to enforce least-privilege access. The workshop `tailnet` has fully open ACLs for simplicity; in production you would give each user and each agent only the access they need.

## Wrapping Up

In this exercise you:

- Pointed the OpenAI client at Aperture so LLM calls flow over the `tailnet` with Tailscale identity as the auth layer, not a client-side API key
- Ran a tool-calling Workflow where the LLM chose a single tool
- Turned that Workflow into an agentic loop where the LLM keeps calling tools until it has enough information to answer
- Used Temporal's dynamic activities to dispatch whichever tool the LLM chose on each iteration
- Watched Temporal record every LLM call and tool result as part of the Workflow history
- Explored the Aperture UI to see per-request logs, tool calls, and cost attribution for every LLM call the agent made

In the final exercise you'll combine the `tsnet` pattern from Exercise 2 with the Aperture pattern you just used, in a single Go service. A metrics watcher that scrapes a `tailnet`-only endpoint, asks Claude for a health summary, and runs on a Temporal Schedule.
