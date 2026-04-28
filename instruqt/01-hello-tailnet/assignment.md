---
slug: hello-tailnet
id: 9cyw4e8rhxrj
type: challenge
title: Hello Tailnet
teaser: Run your first Workflow on the shared Temporal server through the Tailscale
  network
notes:
- type: text
  contents: |-
    # Welcome to the Workshop!

    A Temporal dev server is running on a remote VPS, exposed to this
    Tailscale network via **temporal-ts-net**. Your Exercise
    Environment has the Tailscale client installed but is not on
    the `tailnet` yet; you'll join it in the first step.

    In this challenge you'll join the `tailnet`, configure the Temporal
    environment, run a geo-IP Workflow, and see it in the
    shared Temporal Web UI.
tabs:
- id: yq44efumcgda
  title: Code Editor
  type: code
  hostname: workshop
  path: /root/workshop
- id: jntvq5itaemz
  title: Worker
  type: terminal
  hostname: workshop
  workdir: /root/workshop
- id: q3jgtexmlhzg
  title: Starter
  type: terminal
  hostname: workshop
  workdir: /root/workshop
- id: j2mxffaqlfnf
  title: Temporal UI
  type: service
  hostname: workshop
  port: 8233
difficulty: basic
timelimit: 900
enhanced_loading: null
---

# Exercise 1: Hello `tailnet`

Your first Workflow on the shared Temporal server, accessed through the Tailscale network.

## Background

A Temporal dev server is running on a remote VPS, exposed to this Tailscale network via [temporal-ts-net](https://github.com/temporal-community/temporal-ts-net). **Once you join the tailnet in Step 1**, you'll be able to reach it at `temporal-dev:7233` (gRPC) and via the **Temporal UI** tab. (The UI tab may show a connection error until Step 1 completes; that's expected.)

The Workflow you'll run gets your exercise environment's public IP address, then geolocates it. This exercise is the foundation the rest of the workshop builds on: every later exercise assumes Workers and the Temporal Server can reach each other over the `tailnet` rather than the public internet.

## Environment

All code for this exercise lives in `exercises/01_hello_tailnet/`. Inside that directory:

- **`practice/`** is where you do your work. Each file has one or more **TODO** comments pointing at the change you need to make.
- **`solution/`** contains the finished version of every file. If you get stuck or want to double-check your work, compare against the matching file in `solution/`. Don't run from `solution/`, run from `practice/`.

## Step 1: Join the `tailnet`

Your Exercise Environment has the Tailscale client and an auth key available as `$TS_AUTHKEY`.

In the [button label="Worker" background="#444CE7"](tab-1) terminal, bring Tailscale up:

```bash,run
tailscale up --auth-key="$TS_AUTHKEY" --hostname="${WORKSHOP_USER_ID}-env"
```

Once the command has completed, confirm you're connected and can see the shared Temporal server:

```bash,run
tailscale status
```

You should see all of the devices that are on the tailnet, including the Temporal Server `temporal-dev` in the list.

## Step 2: Verify the Temporal Config File

Temporal CLI and SDKs support configuring a Temporal Client using environment variables and TOML configuration files, rather than setting connection options programmatically in your code. This decouples connection settings from application logic, making it easier to manage different environments such as development, staging, and production without code changes.

This has already been set up for you in this environment. To verify, open `temporal.toml` in the [button label="Code Editor" background="#444CE7"](tab-0) tab. It's already in place in the workshop directory and the SDK is pointed at it via `TEMPORAL_CONFIG_FILE`. You should see two profiles: `default` (localhost) and `tailnet` (pointing at `temporal-dev:7233`).

## Step 3: Configure Your Application to Connect to the `tailnet` Profile

Both `worker.py` and `starter.py` currently load the `default` profile, which points at localhost. You need to point them at the `tailnet` profile so they connect to the shared Temporal server. This is a one-line change in each file.

### 3a. Edit `exercises/01_hello_tailnet/practice/worker.py`

Find the **TODO** in `worker.py` and pass `profile="tailnet"` to `load_client_connect_config`.

**Before:**

```python
# TODO: Load the "tailnet" profile
config = ClientConfig.load_client_connect_config()
```

**After:**

```python
config = ClientConfig.load_client_connect_config(profile="tailnet")
```

### 3b. Edit `exercises/01_hello_tailnet/practice/starter.py`

The same change in `starter.py`:

**Before:**

```python
# TODO: Load the "tailnet" profile
config = ClientConfig.load_client_connect_config()
```

**After:**

```python
config = ClientConfig.load_client_connect_config(profile="tailnet")
```

Now both the worker and starter read the `tailnet` profile from `temporal.toml` and connect to the shared Temporal server.

## Step 4: Add your name to the Workflow ID

Workflow IDs must be unique on a Temporal Server, and prefixing yours with your name lets you find your run among everyone else's in the shared Temporal UI. `USER_ID` is already set from the sign-up form as `WORKSHOP_USER_ID` and wired into `starter.py`.

In `starter.py`, find the **TODO** on the `execute_workflow` call and prefix the Workflow ID with `{USER_ID}-`.

**Before:**

```python
id=f"geo-ip-{uuid.uuid4()}",
```

**After:**

```python
id=f"{USER_ID}-geo-ip-{uuid.uuid4()}",
```

## Step 5: Start the Worker

Now you are ready to run your Workflow. First, start the Worker. This is the process that will execute your Temporal application.

In the [button label="Worker" background="#444CE7"](tab-1) terminal:

```bash,run
cd exercises/01_hello_tailnet/practice
uv run worker.py
```

You should see: `INFO:root:Connecting to Temporal at temporal-dev:7233`

Once it has connected you will see output similar to: `INFO:root:Starting worker on task queue: your-name-hello-tailnet`

## Step 6: Run the Workflow

Once the Worker has started you are ready to run your Workflow.

In the [button label="Starter" background="#444CE7"](tab-2) terminal:

```bash,run
cd exercises/01_hello_tailnet/practice
uv run starter.py
```

You should see the public IP address and location of the server where your exercise environment is hosted.

Sample output:

```bash,nocopy
Your IP address: 257.257.257.257
Your location:   Alderaan, Core Worlds
```

## Step 7: Check the Temporal UI

Click the [button label="Temporal UI" background="#444CE7"](tab-3) tab and find your Workflows by searching for your user ID. You should see your Workflow, along with all other attendees in the workshop!

> **Note:** The **Temporal UI** tab was loaded before your Exercise Environment joined the `tailnet`, so its first render is stale. Click the refresh button at the top of the tab to reload it. If it still shows a connection error, wait a few seconds and refresh again.

**What happened**

Your Worker is running in your exercise environment hosted in a GCP region. The Temporal Server is running on a VPS in DigitalOcean's San Francisco region. Tailscale provided secure access for the Worker to communicate with the server without exposing the server to the public internet.

## Step 8: Route egress through an exit node

Tailscale also allows for setting an exit node on your `tailnet`, acting as a full tunnel VPN.

The `tailnet` has a shared exit node called `nyc3-exit-node`. Route your public internet traffic through it and re-run the Workflow to see your IP and location change.

In the [button label="Starter" background="#444CE7"](tab-2) terminal, set your environment to use `nyc3-exit-node` as the exit node:

```bash,run
tailscale set --exit-node=nyc3-exit-node
```

Then re-run the Workflow:

```bash,run
uv run starter.py
```

Your IP address and location should now come from New York. The Workflow still reaches `temporal-dev` directly over the `tailnet`, because exit nodes only affect egress to the public internet, not `tailnet`-internal traffic.

Unset your exit node before moving on:

```bash,run
tailscale set --exit-node=
```

## Wrapping Up

In this exercise you:

- Joined a `tailnet` and connected to a Temporal Server that is not reachable from the public internet
- Configured a Temporal Client with a TOML profile instead of hard-coding connection settings
- Ran a Workflow from an environment in one cloud region against a Temporal Server in another, with Tailscale as the only network path between them
- Routed your public-internet egress through a shared `tailnet` exit node and saw your geo-IP change while the Workflow's connection to `temporal-dev` was unaffected

The pattern you just set up, a Worker on one side of a `tailnet` talking to a Temporal Server on the other, is the foundation every remaining exercise builds on. In the next exercise you'll take the network one layer deeper by embedding Tailscale directly inside a Go Worker with `tsnet`, so the Worker becomes its own node on the `tailnet` rather than piggybacking on the host's Tailscale client.
