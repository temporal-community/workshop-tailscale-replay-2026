---
slug: hello-tailnet
id: 9cyw4e8rhxrj
type: challenge
title: Hello Tailnet
teaser: Run your first workflow on the shared Temporal server through the Tailscale
  network
notes:
- type: text
  contents: |-
    # Welcome to the Workshop!

    A Temporal dev server is running on a remote VPS, exposed to this
    Tailscale network via **temporal-ts-net**. Your Exercise
    Environment has the Tailscale client installed but is not on
    the tailnet yet; you'll join it in the first step.

    In this challenge you'll join the tailnet, configure the Temporal
    environment, run a simple geo-IP workflow, and see it in the
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

# Exercise 1: Hello Tailnet

Your first workflow on the shared Temporal server — accessed through the Tailscale network.

## Background

A Temporal dev server is running on a remote VPS, exposed to this Tailscale network via [temporal-ts-net](https://github.com/temporal-community/temporal-ts-net). You can reach it at `temporal-dev:7233` (gRPC) and the **Temporal UI** tab above.

The workflow you'll run is simple: it gets your machine's public IP address, then geolocates it.

## Step 1: Join the tailnet

Your Exercise Environment has a Tailscale auth key available as `$TS_AUTHKEY`. In the **Worker** terminal, bring Tailscale up:

```bash
tailscale up --auth-key="$TS_AUTHKEY" --hostname="${WORKSHOP_USER_ID}-env"
```

Confirm you're connected and can see the shared Temporal server:

```bash
tailscale status
```

You should see `temporal-dev` in the list.

## Step 2: Open the Temporal config file

Open `temporal.toml` in the **Code Editor** tab. It's already in place in the workshop directory and the SDK is pointed at it via `TEMPORAL_CONFIG_FILE`. You should see two profiles: `default` (localhost) and `tailnet` (pointing at `temporal-dev:7233`).

## Step 3: Point the SDK at the tailnet profile

Open `exercises/01_hello_tailnet/practice/worker.py` and `starter.py` in the Code Editor. Each one currently loads the default profile, which points at localhost. Find the **TODO** in each file and change:

```python
config = ClientConfig.load_client_connect_config()
```

to:

```python
config = ClientConfig.load_client_connect_config(profile="tailnet")
```

Now the worker and starter will read the `tailnet` profile from `temporal.toml` and connect to the shared Temporal server.

## Step 4: Add your name to the workflow ID

Still in `starter.py`, find the **TODO** on the `execute_workflow` call and add your `USER_ID` to the workflow ID so you can find it in the shared Temporal UI:

```python
id=f"{USER_ID}-geo-ip-{uuid.uuid4()}",
```

## Step 5: Start the worker

In the **Worker** terminal:

```bash
cd exercises/01_hello_tailnet/practice
uv run worker.py
```

You should see: `Connecting to Temporal at temporal-dev:7233`

## Step 6: Run the workflow

In the **Starter** terminal:

```bash
cd exercises/01_hello_tailnet/practice
uv run starter.py
```

You should see your public IP address and location.

## Step 7: Route egress through an exit node

The tailnet has a shared exit node called `nyc3-exit-node`. Route your public internet traffic through it and re-run the workflow to see your IP and location change.

In the **Worker** terminal (the worker can keep running — this only affects outbound internet traffic):

```bash
tailscale down
tailscale up --auth-key="$TS_AUTHKEY" --exit-node=nyc3-exit-node
```

Then in the **Starter** terminal, run the workflow again:

```bash
uv run starter.py
```

Your IP address and location should now come from New York. The workflow still reaches `temporal-dev` directly over the tailnet — exit nodes only affect egress to the public internet, not tailnet-internal traffic.

## Step 8: Check the Temporal UI

Click the **Temporal UI** tab and find your workflows by searching for your user ID. You should see both runs: one with your original location, one from NYC.

You're running workflows on a shared Temporal server, accessible only through your Tailscale network!
