# Running without Instruqt

For instructors delivering a cohort without Instruqt: a meetup, an internal training, a conference breakout, or a brown-bag. Attendees bring their own laptops; you provide the shared Temporal server, the tailnet, and (ideally) an Aperture endpoint.

If you're just trying this on your own machine, see [Try it locally](try-it-locally.md).
If you're running on Instruqt, see [Running with Instruqt](run-with-instruqt.md).

## What you're signing up for

Three pieces you own:

1. **A VPS** - running `temporal-ts-net`, provisioned once before the event and torn down after. See [Infrastructure](infrastructure.md) for the details.
2. **A tailnet** - with reusable auth keys for attendees to join.
3. **An LLM endpoint** - ideally an Aperture deployment holding a shared OpenAI key. Fall back to per-attendee keys if Aperture isn't available.

Plus the ~15 minutes of onboarding friction as everyone installs Tailscale and joins the tailnet. Budget time for that at the start of the session.

## Attendee prerequisites

Send this to attendees a day or two before the workshop:

- Python 3.13 and `uv` ([installer](https://docs.astral.sh/uv/))
- Go 1.26+ ([go.dev/dl](https://go.dev/dl/)) - needed for Exercise 2 and Exercise 4
- The Tailscale client ([tailscale.com/download](https://tailscale.com/download))
- This repo cloned: `git clone https://github.com/temporal-community/workshop-tailscale-replay-2026`
- Dependencies installed: `cd workshop-tailscale-replay-2026 && uv sync`

## Joining the tailnet

At the start of the workshop, hand out the attendee auth key (via chat, slide, or QR code) and have everyone run:

```shell
tailscale up --authkey=tskey-auth-<key> --hostname=<their-name>
```

Verify they can reach the shared server:

```shell
tailscale ping temporal-dev
tailscale status | grep temporal-dev
```

If someone's `tailscale ping` fails, usually it's either a typo in the auth key or a corporate firewall blocking outbound UDP 41641. Have them fall back to [Tailscale's DERP relay](https://tailscale.com/kb/1257/connection-types) mode. No config change needed, just slower.

## Temporal profile

Each attendee needs `~/.config/temporalio/temporal.toml` pointing at your shared server:

```toml
[profile.tailnet]
address = "temporal-dev:7233"
namespace = "default"

[profile.tailnet.env]
APERTURE_URL = "http://ai"
```

No API key env var. Aperture injects the real credential on the server side, so attendee machines never hold it — the exercises pass `api_key=""` to the SDK and Aperture authenticates via Tailscale WhoIs.

Then `export TEMPORAL_PROFILE=tailnet` in every terminal. Exercises run exactly as they would in Instruqt:

```shell
cd exercises/01_hello_tailnet/practice
uv run worker.py
uv run starter.py
```

Workflows show up in the shared UI at `http://temporal-dev:8233`. The whole room sees them, which is the point.

## Aperture, or what to do if you don't have it

The workshop's rate-limit demo (Exercise 3's "everyone fire at once" moment) relies on Aperture enforcing per-identity quotas. If you have Aperture, point `APERTURE_URL` at it and the demo works as designed.

If you don't have Aperture, you have three options:

1. **Each attendee brings their own OpenAI key.** Attendees set `APERTURE_URL=https://api.openai.com` in their `temporal.toml` and change the exercises' `api_key=""` to `api_key=os.getenv("OPENAI_API_KEY")` (with `OPENAI_API_KEY=sk-...` exported). You skip the rate-limit demo. Simplest; requires every attendee to have an OpenAI account.
2. **Shared-key proxy.** Stand up a small reverse proxy on the VPS that injects a single OpenAI key and enforces rate limits yourself (for example, with nginx and `limit_req`). Skips identity-awareness but preserves the shared-key story.
3. **Skip the LLM content.** Cover Exercises 1 and 2; drop 3 and 4. Not a great tradeoff. The agent content is the payoff.

Pick (1) if your audience are pro devs with their own accounts. Pick (2) if you want to preserve the shared-key teaching moment.

## Pre-event checklist

- [ ] VPS provisioned and `temporal-ts-net` running as a service (see [Infrastructure](infrastructure.md))
- [ ] `temporal-dev:7233` reachable from a Tailscale-connected machine
- [ ] `temporal-dev:8233` (Web UI) loads in a browser
- [ ] Attendee reusable auth key generated and working (join a test machine, confirm it appears in `tailscale status`)
- [ ] Aperture endpoint working, or fallback plan finalized (option 1, 2, or 3 above)
- [ ] Attendees notified of prerequisites and given the repo link
- [ ] You have walked through all four exercises end-to-end on a fresh machine since your last rehearsal

## Day-of checklist

- [ ] VPS healthy (`tailscale ping temporal-dev` from your machine)
- [ ] Temporal UI loads
- [ ] LLM path working: run one agent end-to-end before the first attendee connects
- [ ] Auth key queued up for sharing (QR code printed, or paste-ready in chat)
- [ ] Slides open

## Slides

The deck lives at `slides/slides.md` and runs under Slidev. See [Present the slides](slides.md) for installing pnpm, running the deck, presenter mode, and exporting to PDF.

## Backup plan

If the shared server fails mid-workshop:

1. Have attendees run `temporal server start-dev` locally.
2. They edit `temporal.toml`: switch the `default` profile's address to `localhost:7233` and `export TEMPORAL_PROFILE=default`.
3. LLM calls keep working independently (they still go through Aperture via the tailnet, assuming Aperture didn't also fail).
4. They lose the shared-UI "see everyone's workflows" experience, but every exercise still completes.

Rehearse this beforehand. Switching an attendee over live is frustrating if it's the first time you've tried.

## Post-workshop

1. Rotate or revoke the attendee auth key in the Tailscale admin console. Otherwise it's a long-lived credential sitting in whoever's terminal history.
2. If the VPS was purpose-built, destroy it once you've grabbed anything you want to keep from `/var/lib/temporal/workshop.db`.
3. Share the [workshop repo](https://github.com/temporal-community/workshop-tailscale-replay-2026) and [Try it locally](try-it-locally.md) with attendees who want to keep exploring.
