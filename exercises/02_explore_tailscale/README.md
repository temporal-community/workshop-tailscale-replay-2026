# Exercise 2: Exploring Your Tailscale Network

Now that you've run a workflow through the tailnet, let's understand the network you're on.

## Goal

Use Tailscale CLI tools to discover what's on the network, understand how access control works, see how Aperture fits into the architecture, and run the same geo-IP workflow from Exercise 1 in Go — but this time the worker joins the tailnet itself via [tsnet](https://pkg.go.dev/tailscale.com/tsnet) instead of riding on the VM's Tailscale client.

## Part 1: Discover Your Network

### See all machines on the tailnet

```bash
tailscale status
```

You should see:
- **Your machine** — the Instruqt VM you're working on
- **temporal-dev** — the VPS running the shared Temporal dev server
- **Other attendee machines** — everyone else in the workshop
- **Aperture endpoint** — the API gateway that will proxy your LLM calls

### Ping the Temporal server

```bash
tailscale ping temporal-dev
```

Notice the latency. This is a direct, encrypted WireGuard connection — no relay servers, no VPN concentrators.

### Check what's accessible

Try reaching the Temporal gRPC port and Web UI:

```bash
# gRPC port (should succeed)
nc -zv temporal-dev 7233

# Web UI (should succeed — you already used this in Exercise 1)
curl -s -o /dev/null -w "%{http_code}" http://temporal-dev:8233
```

### Check your Tailscale identity

```bash
tailscale whois $(tailscale ip -4)
```

This shows your identity on the tailnet, and that can be used by devices and applications on the tailnet to authenticate you.

## Part 2: Understanding Aperture

Aperture by Tailscale is a centralized AI gateway that secures, monitors, and routes LLM requests across your organization. Aperture uses Tailscale's identity layer (WhoIs) to automatically authenticate users, eliminating the need to distribute API keys. It routes requests to upstream LLM providers such as OpenAI, Anthropic, and Google without requiring changes to existing tools or workflows.

There are four core mechanisms: identity and authentication, request routing by model, telemetry capture and session logging. Together, these enable auditing, cost awareness, and operational insight for letting teams safely adopt LLMs.

**Identity and Authentication**: Traditional API proxies require clients to authenticate with tokens or API keys. Aperture eliminates this step by using Tailscale's identity layer where user identity is passed through headers from the client devices running Tailscale to Aperture that is also tailscale-aware and able to validate these headers and thus user identity. User identity itself is autheticated via an IDP and establishes a trusted identity within the tailnet.

**Request Routing by Model**: When a request arrives, Aperture extracts the model name from the request body (for example, `claude-sonnet-4-6` or `gpt-4o`). The proxy looks up which provider serves that model and forwards the request to that provider's API endpoint, injecting the correct authentication headers. From the client's perspective, the proxy appears as if it were the LLM provider itself. Clients connect to the proxy URL and send standard API requests. The proxy handles the routing transparently.

**Telemetry Capture**: The capture system records everything needed to reconstruct and analyze each LLM interaction including request and response headers and body, token counts by type (input,output,cached,reasoning), model name, request duration, tool use count, session context and other metadata. Related requests are grouped into sessions for easier analysis.

**Extensibility using Webhooks**: Aperture provides a webhook mechanism that allows you to extend its functionality. You can use webhooks to integrate with external systems, perform custom logic, or trigger actions based on events. Some of the existing partners that have integrated with Aperture include [Oso](https://www.osohq.com/), [Cerbos](https://cerbos.dev/), and [Cribl](https://www.cribl.io/)

### Aperture endpoint

The Aperture proxy is available on the workshop tailnet at:

```
http://ai
```

This endpoint accepts OpenAI-compatible requests and forwards them to OpenAI with the shared API key. Your Tailscale identity is used to enforce per-user rate limits — no API key lives on your machine.

### Try it out

```bash
curl http://ai/v1/models
```

## Part 3: Hello World in Go with tsnet

In Exercise 1 your Python worker reached the Temporal server through the VM's system Tailscale client — the VM joined the tailnet, and the worker inherited that connectivity. There's another way: [tsnet](https://pkg.go.dev/tailscale.com/tsnet) lets a program join the tailnet by itself, as its own node, with no system Tailscale install required.

You're going to run the same geo-IP workflow from Exercise 1, but in Go, with the worker joining the tailnet as its own node.

### Why tsnet?

- **Self-contained binaries.** Ship a single Go binary that joins the tailnet — no apt install, no daemon, no separate `tailscale up`.
- **Per-process identity.** The worker is its own node on the tailnet, with its own hostname, its own ACL-eligible identity. The VM it runs on doesn't need to be on the tailnet at all.
- **Same tailnet, different door.** The `temporal-dev` server doesn't care how you connect — system client or tsnet, the traffic is the same encrypted WireGuard.

*Note: Aperture itself is a tsnet application, so it can join the tailnet just like any other workload running the tailscale client, as is the Temporal Dev server!

### How the code is laid out

```
exercises/02_explore_tailscale/go-hello-tsnet/practice/
├── go.mod
├── activities.go      # GetIP + GetLocationInfo — identical to Ex 1 in Go
├── shared.go          # WorkflowInput/Output types
├── workflow.go        # GetAddressFromIP workflow
└── main.go            # tsnet setup + Temporal worker/starter (TODOs here)
```

Activities and workflow are done. You'll fill in `main.go` to join the tailnet via tsnet and dial Temporal through it.

### Step 1: Move to the practice directory

```bash
cd exercises/02_explore_tailscale/go-hello-tsnet/practice
go mod download
```

### Step 2: Complete TODO 1 — Start a tsnet node

Open `main.go`. Find TODO 1. Create a `tsnet.Server` with your own hostname and start it:

```go
nodeName := fmt.Sprintf("%s-ex2-go-%s-%s", userID, mode, randSuffix)
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

The `Dir` field holds the tsnet state (node key, machine key). On first run it uses `TS_AUTHKEY` to register; on subsequent runs it reuses the stored identity. The random 5-char suffix is generated once and persisted via the state dir — so two attendees with the same `WORKSHOP_USER_ID` land on different hostnames on the tailnet.

### Step 3: Complete TODO 2 — Dial Temporal through tsnet

The Temporal SDK opens a gRPC connection to `temporal-dev:7233`. To route that through the tailnet, add a `grpc.WithContextDialer` option that calls `tsNode.Dial`.

Inside `dialTemporal`, replace the `// TODO 2` comment in the `dialOptions` slice with the `grpc.WithContextDialer(...)` call. The slice should end up looking like:

```go
dialOptions := []grpc.DialOption{
    grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
        return tsNode.Dial(ctx, "tcp", addr)
    }),
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithKeepaliveParams(keepalive.ClientParameters{
        Time:                30 * time.Second,
        Timeout:             10 * time.Second,
        PermitWithoutStream: true,
    }),
}
```

Every byte the SDK sends now goes through `tsNode.Dial`, routed over the tailnet.

### Step 4: Run the worker

```bash
TS_AUTHKEY=tskey-auth-<your-key> \
WORKSHOP_USER_ID=$WORKSHOP_USER_ID \
go run . worker
```

First run takes ~5 seconds to join the tailnet. You should see:

```
joined tailnet as <your-user-id>-ex2-go-worker-<5 random chars>
Starting Go worker on task queue: <your-user-id>-hello-tsnet
```

### Step 5: Confirm the worker is on the tailnet

In another terminal:

```bash
tailscale status | grep $WORKSHOP_USER_ID-ex2-go-worker
```

You should see your Go worker listed as its own node — distinct from the VM itself. That's tsnet.

### Step 6: Run the workflow

```bash
TS_AUTHKEY=tskey-auth-<your-key> \
WORKSHOP_USER_ID=$WORKSHOP_USER_ID \
go run . starter
```

The starter uses the same tsnet pattern to dial Temporal, fires `GetAddressFromIP`, and prints the result. Find it in the Temporal UI at `http://temporal-dev:8233`.

### What changed from Exercise 1?

- The workflow and activities are the same shape — same two activities (`GetIP`, `GetLocationInfo`), same pattern.
- The language is Go instead of Python.
- The worker's network identity is different: it's `<user>-ex2-go-worker-<suffix>` on the tailnet, not the VM's hostname.
- No system Tailscale client is involved in this worker's path.

This is the pattern Exercise 4 (stretch) uses to run a Go agent with Aperture. Same `tsNode.Dial` trick, different workflow.

## What You've Learned

- How to discover machines on a Tailscale network
- That Tailscale provides identity for free — no extra auth layer needed
- How Aperture uses that identity to proxy and rate-limit API calls
- How to embed Tailscale directly into a Go worker with tsnet
- That the same Temporal workflow runs identically whether the worker joins the tailnet via the system client (Ex 1) or via tsnet (Ex 2 Part 3)
- The full network path: Your VM → Tailscale → Temporal Server / Aperture → OpenAI

Next up: you'll use this Aperture endpoint to power an AI agent workflow.
