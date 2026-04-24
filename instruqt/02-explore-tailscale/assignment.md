---
slug: explore-tailscale
id: jysnh4r6f6cu
type: challenge
title: Explore Your Tailscale Network
teaser: Explore the tailnet and run a Go Worker that joins it via tsnet
notes:
- type: text
  contents: |-
    # Exploring the Network

    Now that you've proven the `tailnet` works, let's look under the hood.
    You'll discover what's on the `tailnet` and then run a Go Worker
    that joins the `tailnet` as its own node via `tsnet`, instead of
    riding on the Exercise Environment's Tailscale client.
tabs:
- id: 3rfluijsiy9t
  title: Worker
  type: terminal
  hostname: workshop
  workdir: /root/workshop
- id: sx3uajfowyxv
  title: Starter
  type: terminal
  hostname: workshop
  workdir: /root/workshop
- id: fi2edgfrbxys
  title: Code Editor
  type: code
  hostname: workshop
  path: /root/workshop/exercises/02_explore_tailscale
- id: le7rxdsm3qqx
  title: Temporal UI
  type: service
  hostname: workshop
  port: 8233
difficulty: basic
timelimit: 1200
enhanced_loading: null
---

# Exercise 2: Exploring Your Tailscale Network

Now that you've run a Workflow through the `tailnet`, let's look at the network itself. You'll discover what's on the `tailnet`, see how Tailscale identity works, and run a Go Worker that joins the `tailnet` as its own node via [`tsnet`](https://pkg.go.dev/tailscale.com/tsnet) instead of riding on the Exercise Environment's system client.

## Background

In Exercise 1 your Python Worker reached Temporal through the Exercise Environment's Tailscale client. The environment was on the `tailnet`, and the Worker inherited that connectivity. `tsnet` is the library inside Tailscale that lets a process join the `tailnet` directly, as its own node, with no system-wide install. Every Tailscale binary uses it under the hood, and you can embed it in any Go program. That's what `temporal-ts-net` does to put the Temporal dev server on the `tailnet`, and it's the pattern you'll use in Exercise 4 for the metrics watcher.

## Environment

Steps 1 through 3 poke at the `tailnet` from the **Worker** terminal and don't touch code. Step 4 takes the system Tailscale client offline so the later `tsnet` work has to stand on its own. Starting in Step 5, the Go code for this exercise lives under `go-hello-tsnet/` in the **Code Editor** tab. Inside that directory:

- **`practice/`** is where you do your work. Each file has one or more **TODO** comments pointing at the change you need to make.
- **`solution/`** contains the finished version of every file. If you get stuck or want to double-check your work, compare against the matching file in `solution/`. Don't run from `solution/`, run from `practice/`.

Go and the `tsnet` library are already installed, and the module cache has been warmed up by the workshop setup, so `go run` builds quickly without fetching from the network.

> **Verify you're on the `tailnet`**
>
> Run the following command:
> ```bash
> tailscale status
> ```
>
> If you see **Logged Out** then you need to reauthenticate to the `tailnet`
>
> Run the following command to authenticate to the `tailnet`
> ```bash
> tailscale up --auth-key="$TS_AUTHKEY" --hostname="${WORKSHOP_USER_ID}-env"
> ```

## Step 1: See what's on the `tailnet`

Before you add another Worker, it's worth looking at who else is already on the `tailnet`. `tailscale status` lists every node your machine can see, along with its IP, hostname, and connection path.

In the **Worker** terminal:

```bash
tailscale status
```

You should see:

- **Your Exercise Environment** (`workshop-<something>`) with a `100.x.y.z` `tailnet` IP
- **`temporal-dev`**, the VPS running the shared Temporal dev server
- **Other attendee machines**, everyone else in the workshop, each with their own hostname

## Step 2: Ping the Temporal Server

`tailscale status` tells you what's on the `tailnet`. `tailscale ping` tells you *how* you reach any given node, whether the first packets go through a Tailscale relay (DERP) or straight over a direct encrypted WireGuard path.

Still in the **Worker** terminal:

```bash
tailscale ping temporal-dev
```

You should see output similar to:

```output
pong from temporal-dev (100.109.42.22) via 167.71.156.227:35753 in 40ms
```

The first line may say `pong from temporal-dev (100.109.42.22) via DERP(...)`, meaning it was relayed through Tailscale's infrastructure. Subsequent lines (like the sample above) report `via <public-IP>:<port>`, a direct encrypted WireGuard path with no relay. Once the direct path is established, every packet to `temporal-dev` flows over it.

## Step 3: Check your Tailscale identity

Every node on the `tailnet` has an identity: a machine name, tags, and the user that owns the node. Services on the `tailnet` use this identity to authorize or rate-limit you without any API keys on your side. In Exercise 3 you'll see Aperture use exactly this mechanism to attach your name to every LLM call.

Still in the **Worker** terminal:

```bash
tailscale whois $(tailscale ip -4)
```

You should see the machine name, tags, and the Tailscale user associated with this node.

## Step 4: Take the system Tailscale client offline

Steps 1 through 3 leaned on the system `tailscale` binary to explore the `tailnet`. The Go Worker you're about to build joins the `tailnet` on its own via `tsnet`, so it doesn't need the binary at all. Prove that by taking it down now, before the Worker ever runs.

In the **Worker** terminal:

```bash
tailscale down
```

Verify the client is stopped:

```bash
tailscale status
```

You should see `Stopped` (or `Logged out`) instead of the `tailnet` nodes from Step 1. Your Exercise Environment is no longer on the `tailnet`. The Go Worker is about to put itself there on its own.

## Step 5: Configure the `tsnet.Server`

The `tsnet.Server` struct is how a Go program declares that it wants to be a `tailnet` node. You give it a hostname, a state directory, and an auth key, and once you call `Start()` your process has its own `tailnet` IP.

Open `go-hello-tsnet/practice/main.go` in the **Code Editor** tab. Inside `startTsnet`, find the **TODO** in the `tsnet.Server` literal and set each field using the values already in scope:

- Set **`Hostname`** to `nodeName`. This is the name your node will have in `tailscale status`. The `resolveNodeName` helper just above the TODO has already computed it as `<userID>-ex2-go-<mode>-<5 random chars>`.
- Set **`Dir`** to `filepath.Join(configDir, "workshop-tsnet", nodeName)`. This is where `tsnet` persists its node and machine keys so later runs reuse the identity.
- Set **`AuthKey`** to `os.Getenv("TS_AUTHKEY")`. This is consumed once on first run to register the node on the `tailnet`.

The finished struct literal should look like this:

```go
tsNode := &tsnet.Server{
    Hostname: nodeName,
    Dir:      filepath.Join(configDir, "workshop-tsnet", nodeName),
    AuthKey:  os.Getenv("TS_AUTHKEY"),
}
```

The rest of the function (`tsNode.Start()`, `tsNode.Up(upCtx)`, the log line) is already in place.

## Step 6: Dial Temporal through `tsnet`

The `tsnet.Server` you just configured gives your Worker its own node on the `tailnet`. Now you need to tell the Temporal Go SDK to use it, so the gRPC connection to `temporal-dev:7233` flows through `tsnet` instead of the system network stack.

Still in `main.go`, find the **TODO** inside `dialTemporal`. Add a `grpc.WithContextDialer` entry at the top of the `dialOptions` slice. The rest of the slice is already in place:

```go
grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
    return tsNode.Dial(ctx, "tcp", addr)
}),
```

Every byte the SDK sends now flows through `tsNode.Dial`, which routes over the `tailnet`.

## Step 7: Start the Go Worker

With both TODOs filled in, the Worker is ready to join the `tailnet` and connect to Temporal on its own.

In the **Worker** terminal:

```bash
cd exercises/02_explore_tailscale/go-hello-tsnet/practice
go run . worker
```

The first run takes about five seconds while `tsnet` registers the node. You should see output similar to:

```output
joined tailnet as <your-user-id>-ex2-go-worker-<5 random chars>
connected to temporal at temporal-dev:7233 via tsnet
Starting Go worker on task queue: <your-user-id>-hello-tsnet
```

## Step 8: Confirm the Go Worker is on the `tailnet`

The Go Worker should now appear as its own node on the `tailnet`, even with the system Tailscale client still offline from Step 4. To verify from the command line, bring the system client back up and check `tailscale status`.

In the **Starter** terminal:

```bash
tailscale up --auth-key="$TS_AUTHKEY" --hostname="${WORKSHOP_USER_ID}-env"
tailscale status | grep -- '-ex2-go-worker'
```

You should see a new row, `<your-user-id>-ex2-go-worker-<suffix>`, separate from the Exercise Environment itself. The Worker joined the `tailnet` on its own via `tsnet` while the system client was offline. That's the whole point of Step 4.

## Step 9: Run the Workflow

Now trigger the same geo-IP Workflow from Exercise 1. This time the activities execute on the Go `tsnet` Worker you just started.

Still in the **Starter** terminal:

```bash
cd exercises/02_explore_tailscale/go-hello-tsnet/practice
go run . starter
```

You should see your public IP address and location printed, same as Exercise 1, but this time the Worker that executed the activities was the Go `tsnet` Worker, not the Python Worker from Exercise 1.

## Step 10: Check the Temporal UI

Click the **Temporal UI** tab and find your `<your-user-id>-hello-tsnet` Workflow. Click into it and look at the worker info on each activity. The task queue is `<your-user-id>-hello-tsnet`, and the worker identity reflects the Go process rather than the Python one from Exercise 1.

**What happened**

Same Temporal Server, same Workflow, different Worker transport. The Python Worker in Exercise 1 relied on the Exercise Environment's Tailscale client to reach Temporal. The Go Worker you just ran carries its own `tsnet` node inside the process itself, joins the `tailnet` on startup, and dials Temporal through that embedded node.

## Wrapping Up

In this exercise you:

- Inspected the `tailnet` from the command line with `tailscale status`, `tailscale ping`, and `tailscale whois`
- Took the system Tailscale client offline to prove that `tsnet` does not depend on it
- Configured a `tsnet.Server` inside a Go Worker so the Worker itself became a `tailnet` node, not a client of one
- Wired the Temporal Go SDK to dial through `tsnet` with a `grpc.WithContextDialer`
- Ran the same geo-IP Workflow as Exercise 1, but executed by the Go `tsnet` Worker while the system client was down

`tsnet` is the building block the rest of the workshop uses to put Workers and servers anywhere on the `tailnet` with no network wiring. In the next exercise you'll build a Python AI agent whose LLM calls flow through **Aperture**, a `tailnet`-only gateway that uses your Tailscale identity as the auth layer instead of API keys.
