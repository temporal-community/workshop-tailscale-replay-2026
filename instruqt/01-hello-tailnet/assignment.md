---
slug: hello-tailnet
type: challenge
title: "Hello Tailnet"
teaser: Run your first workflow on the shared Temporal server through the Tailscale network
difficulty: basic
timelimit: 900
notes:
  - type: text
    contents: |-
      # Welcome to the Workshop!

      A Temporal dev server is running on a remote VPS, exposed to this
      Tailscale network via **temporal-ts-net**. Your machine is already
      connected to the tailnet.

      In this challenge you'll configure the Temporal environment, run a
      simple geo-IP workflow, and see it in the shared Temporal Web UI.
tabs:
  - title: Terminal
    type: terminal
    hostname: workshop
  - title: Code Editor
    type: code
    hostname: workshop
    path: /root/workshop
  - title: Temporal UI
    type: service
    hostname: workshop
    port: 8233
    url: http://temporal-dev:8233
---

# Exercise 1: Hello Tailnet

Your first workflow on the shared Temporal server — accessed through the Tailscale network.

## Background

A Temporal dev server is running on a remote VPS, exposed to this Tailscale network via [temporal-ts-net](https://github.com/temporal-community/temporal-ts-net). You can reach it at `temporal-dev:7233` (gRPC) and the **Temporal UI** tab above.

The workflow you'll run is simple: it gets your machine's public IP address, then geolocates it.

## Step 1: Verify your environment

```bash
cd /root/workshop
uv run scripts/verify_setup.py
```

All checks should pass. If the Temporal server check fails, verify Tailscale is connected:

```bash
tailscale status
```

## Step 2: Create the Temporal config file

The Temporal SDK reads connection settings from a TOML configuration file. Create it:

```bash
mkdir -p ~/.config/temporalio
cp /root/workshop/temporal.toml.example ~/.config/temporalio/temporal.toml
```

Open `~/.config/temporalio/temporal.toml` in the **Code Editor** tab and verify it has the `tailnet` profile:

```toml
[profile.tailnet]
address = "temporal-dev:7233"
namespace = "default"
```

The `TEMPORAL_PROFILE=tailnet` environment variable is already set, so the SDK will use this profile automatically.

## Step 3: Add your name to the workflow ID

Open `exercises/01_hello_tailnet/practice/starter.py` in the Code Editor.

Find the TODO and add your `USER_ID` to the workflow ID so you can find it in the shared Temporal UI:

```python
id=f"{USER_ID}-geo-ip-{uuid.uuid4()}",
```

## Step 4: Start the worker

```bash
cd /root/workshop/exercises/01_hello_tailnet/practice
uv run worker.py
```

You should see: `Connecting to Temporal at temporal-dev:7233`

## Step 5: Run the workflow

Open a **new terminal** (click the `+` button) and run:

```bash
cd /root/workshop/exercises/01_hello_tailnet/practice
uv run starter.py
```

You should see your public IP address and location.

## Step 6: Check the Temporal UI

Click the **Temporal UI** tab above and find your workflow by searching for your user ID.

You're running workflows on a shared Temporal server, accessible only through your Tailscale network!
