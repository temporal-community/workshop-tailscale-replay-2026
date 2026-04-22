# Exercise 1: Hello Tailnet

Your first workflow on the shared Temporal server, accessed through the Tailscale network.

## Goal

Connect your worker to the shared Temporal dev server running on the tailnet, run a simple geo-IP workflow, and see it in the Temporal Web UI.

## Background

A Temporal dev server is running on a remote VPS, exposed to this Tailscale network via [temporal-ts-net](https://github.com/temporal-community/temporal-ts-net). You can reach it at `temporal-dev:7233` (gRPC) and `http://temporal-dev:8233` (Web UI). No VPN setup, no port forwarding, just Tailscale.

The workflow you'll run is simple: it gets your machine's public IP address, then geolocates it. Two activities, one result.

### How Environment Configuration Works

Instead of hardcoding server addresses in your code, the Temporal SDK reads connection settings from a [TOML configuration file](https://docs.temporal.io/develop/environment-configuration). This lets you switch between environments (local dev, tailnet, cloud) without changing any code. Just switch profiles.

Your job is to create the config file, then point the worker and starter at the `tailnet` profile so they connect to the shared server.

## Instructions

### Step 1: Verify your environment

From the repository root:

```bash
uv run scripts/verify_setup.py
```

All dependency checks should pass. The Temporal connectivity check will fail, that's expected, you haven't configured the connection yet.

### Step 2: Complete TODO 1 - Create the Temporal config file

Create the file `~/.config/temporalio/temporal.toml` (you can copy the template from the repo root):

```bash
mkdir -p ~/.config/temporalio
cp temporal.toml.example ~/.config/temporalio/temporal.toml
```

Open the file and verify it has the `tailnet` profile pointing to the shared server:

```toml
[profile.tailnet]
address = "temporal-dev:7233"
namespace = "default"
```

### Step 3: Complete TODO 2 - Load the tailnet profile

Open `practice/worker.py` and `practice/starter.py`. Each one calls `ClientConfig.load_client_connect_config()` with no arguments, which would load the `default` profile (localhost). Change both calls to load the `tailnet` profile:

```python
config = ClientConfig.load_client_connect_config(profile="tailnet")
```

### Step 4: Complete TODO 3 - Add your name to the workflow ID

Still in `practice/starter.py`. Find the TODO near the `execute_workflow` call and add your `USER_ID` to the workflow ID so you can find it in the shared Temporal UI:

```python
id=f"{USER_ID}-geo-ip-{uuid.uuid4()}",
```

### Step 5: Start the worker

In one terminal:

```bash
cd exercises/01_hello_tailnet/practice
uv run worker.py
```

You should see: `Connecting to Temporal at temporal-dev:7233` followed by `Starting worker on task queue: <your-user-id>-hello-tailnet`

### Step 6: Run the workflow

In a second terminal:

```bash
cd exercises/01_hello_tailnet/practice
uv run starter.py
```

You should see your public IP address and location printed.

### Step 7: Open the Temporal Web UI

Open your browser and go to:

```
http://temporal-dev:8233
```

Find your workflow by searching for your user ID. Click into it to see the execution history. Both activities (get_ip, get_location_info) should show as completed.

You're now running workflows on a shared Temporal server, accessible only through your Tailscale network. The same config file works for the worker, the starter, and the Temporal CLI:

```bash
temporal --profile tailnet workflow list
```

## Stuck?

Check the `solution/` directory for the completed code, and `temporal.toml.example` in the repo root for the config file.
