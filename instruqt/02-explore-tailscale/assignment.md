---
slug: explore-tailscale
id: jysnh4r6f6cu
type: challenge
title: Explore Your Tailscale Network
teaser: Explore the tailnet and run a Go worker that joins it via tsnet
notes:
- type: text
  contents: |-
    # Exploring the Network

    Now that you've proven the tailnet works, let's look under the hood.
    You'll discover what's on the tailnet and then run a Go worker
    that joins the tailnet as its own node via tsnet, instead of
    riding on the Exercise Environment's Tailscale client.
tabs:
- id: fi2edgfrbxys
  title: Code Editor
  type: code
  hostname: workshop
  path: /root/workshop/exercises/02_explore_tailscale
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

Now that you've run a workflow through the tailnet, let's look at the network itself: discover what's on it, see how Tailscale identity works, and run a Go worker that joins the tailnet as its own node via [`tsnet`](https://pkg.go.dev/tailscale.com/tsnet) instead of riding on the Exercise Environment's system client.

## Background

In Exercise 1 your Python worker reached Temporal through the Exercise Environment's Tailscale client: the environment was on the tailnet, and the worker inherited that connectivity. `tsnet` is the library inside Tailscale that lets a process join the tailnet directly, as its own tailnet node, with no system-wide install. Every Tailscale binary uses it under the hood, and you can embed it in any Go program. That's what `temporal-ts-net` does to put the Temporal dev server on the tailnet, and it's the pattern you'll use in Exercise 4 for the metrics watcher.

> **Not on the tailnet?** If you joined late or `tailscale status` shows **Logged out**, run this in the **Worker** terminal first:
>
> ```bash
> tailscale up --auth-key="$TS_AUTHKEY" --hostname="${WORKSHOP_USER_ID}-env"
> ```

## Step 1: See what's on the tailnet

In the **Worker** terminal:

```bash
tailscale status
```

You should see:

- **Your Exercise Environment** (`workshop-<something>`) with a `100.x.y.z` tailnet IP
- **`temporal-dev`** - the VPS running the shared Temporal dev server
- **Other attendee machines** - everyone else in the workshop, each with their own hostname

## Step 2: Ping the Temporal server

```bash
tailscale ping temporal-dev
```

Watch the output. The first line says `pong ... via DERP` (relayed through Tailscale's infrastructure) and subsequent lines say `pong ... via <IP>:<port>` - that's a direct encrypted WireGuard path, no relay. Every packet to `temporal-dev` goes over that path.

## Step 3: Check your Tailscale identity

```bash
tailscale whois $(tailscale ip -4)
```

This prints your tailnet identity - machine name, tags, and the user that owns the node. Every packet you send on the tailnet is attributable to this identity, so services on the tailnet can authorize or rate-limit you without any API keys on your side.

## Step 4: Move to the Go `tsnet` practice directory

The rest of this exercise runs a Go worker that joins the tailnet via `tsnet`, then executes the same geo-IP workflow from Exercise 1.

In the **Worker** terminal:

```bash
cd exercises/02_explore_tailscale/go-hello-tsnet/practice
go mod download
```

## Step 5: Complete **TODO 1** - Start a `tsnet` node

Open `exercises/02_explore_tailscale/go-hello-tsnet/practice/main.go` in the Code Editor. Find **TODO 1**. It's a scaffold for a `tsnet.Server`. Fill it in:

```go
tsNode := &tsnet.Server{
    Hostname: nodeName,
    Dir:      filepath.Join(configDir, "workshop-tsnet", nodeName),
    AuthKey:  os.Getenv("TS_AUTHKEY"),
}
if err := tsNode.Start(); err != nil {
    log.Fatalf("tsnet start: %v", err)
}
defer tsNode.Close()

upCtx, upCancel := context.WithTimeout(context.Background(), 30*time.Second)
defer upCancel()
if _, err := tsNode.Up(upCtx); err != nil {
    log.Fatalf("tsnet up: %v", err)
}
log.Printf("joined tailnet as %s", nodeName)
```

`nodeName` is resolved to `<userID>-ex2-go-<mode>-<5 random chars>` by the `resolveNodeName` helper already in `main.go`. The random suffix is generated once and then reused via the state dir, so your worker and starter each land on a stable hostname that won't collide with another attendee using the same `WORKSHOP_USER_ID`.

`Dir` holds the `tsnet` state (node key, machine key). First run uses `TS_AUTHKEY` to register the node; subsequent runs reuse the stored identity.

## Step 6: Complete **TODO 2** - Dial Temporal through `tsnet`

Still in `main.go`, find **TODO 2** near `dialTemporal`. The Temporal SDK wants to open a gRPC connection to `temporal-dev:7233`. To route that through the tailnet, inject a custom `ContextDialer` that calls `tsNode.Dial`:

```go
c, err := client.Dial(client.Options{
    HostPort: "temporal-dev:7233",
    ConnectionOptions: client.ConnectionOptions{
        DialOptions: []grpc.DialOption{
            grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
                return tsNode.Dial(ctx, "tcp", addr)
            }),
            grpc.WithTransportCredentials(insecure.NewCredentials()),
        },
    },
})
```

Every byte the SDK wants to send now flows through `tsNode.Dial`, which routes over the tailnet.

## Step 7: Start the Go worker

In the **Worker** terminal:

```bash
cd exercises/02_explore_tailscale/go-hello-tsnet/practice
go run . worker
```

First run takes ~5 seconds to register on the tailnet. You should see:

```
joined tailnet as <your-user-id>-ex2-go-worker-<5 random chars>
connected to temporal at temporal-dev:7233 via tsnet
Starting Go worker on task queue: <your-user-id>-hello-tsnet
```

## Step 8: Confirm the Go worker is on the tailnet

In the **Starter** terminal:

```bash
tailscale status | grep -- '-ex2-go-worker'
```

You should see a new row - `<your-user-id>-ex2-go-worker-<suffix>` - separate from the Exercise Environment itself. The worker is its own tailnet node with its own hostname and identity.

## Step 9: Run the workflow

Still in the **Starter** terminal:

```bash
cd exercises/02_explore_tailscale/go-hello-tsnet/practice
go run . starter
```

You should see your public IP address and location printed, same as Exercise 1 - but this time the worker that executed the activities was the Go `tsnet` worker, not a Python worker on the Exercise Environment.

## Step 10: Check the Temporal UI

Click the **Temporal UI** tab and find your `<your-user-id>-hello-tsnet` workflow. Click into it and look at the worker info on each activity: the task queue is `<your-user-id>-hello-tsnet`, and the worker identity reflects the Go process rather than the Python one from Exercise 1.

Same Temporal server, same workflow, different worker transport. `tsnet` is the building block the rest of the workshop uses to put workers and servers anywhere on the tailnet with no network wiring.
