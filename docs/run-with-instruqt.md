# Running with Instruqt

This page covers what the Instruqt VMs need, how attendees join the tailnet, and the pre-event checklists the workshop team uses.

If you're running the workshop on your own hardware, see [Running without Instruqt](run-without-instruqt.md). If you're just exploring the repo on your laptop, see [Try it locally](try-it-locally.md).

## VM requirements

Each attendee VM needs:

- **Python 3.13** - `pyenv` or system package
- **`uv`** - the Python dependency manager the exercises use
- **Go 1.26+** - for the Exercise 2 Go worker and the Exercise 4 metrics watcher
- **Tailscale** - installed and joined to the workshop tailnet
- **Browser access** - to `http://temporal-dev:8233` (Temporal Web UI)

## Lifecycle script

Each Instruqt VM boots with a lifecycle script that:

1. Clones this repo.
2. Writes a `.env` with per-VM values (primarily the attendee's unique ID, used as a workflow-ID prefix).
3. Runs `uv sync` to install Python deps.
4. Runs `go mod download` in each Go exercise directory.

The existing reference script lives at `scripts/setup_env.sh`. Adapt it to your Instruqt track shape.

## Tailscale setup

Each VM needs to join the workshop tailnet. Three options, in order of preference:

1. **Pre-authorized reusable auth key** in the lifecycle script. Simplest, lowest-friction.
2. **Tailscale Connector**, if your Instruqt runtime supports it.
3. **Manual join** as a fallback. Attendees run `tailscale up --authkey ...` from a task step.

Coordinate with whoever owns the tailnet to get a key that's:

- **Reusable** - every VM uses the same key
- **Pre-authorized** - no manual approval in the admin console
- **Tagged** - e.g. `tag:workshop`, so ACLs can gate what attendees reach

## Aperture

The Aperture endpoint the workshop points at (`http://ai`) must be reachable from every attendee VM over the tailnet. That's owned by the Tailscale side of the partnership; confirm the endpoint and model availability with them before the event.

## Pre-event checklist

- [ ] VPS provisioned and `temporal-ts-net` running as a service (see [Infrastructure](infrastructure.md))
- [ ] `temporal-dev:7233` reachable from a Tailscale-connected machine
- [ ] `temporal-dev:8233` (Web UI) loads in a browser
- [ ] Aperture endpoint confirmed working (a test `POST /v1/responses` returns a valid completion)
- [ ] Instruqt track tested end-to-end by someone who hasn't seen it before
- [ ] Workshop repo up to date on `main`
- [ ] `uv run scripts/verify_setup.py` passes from a fresh Instruqt VM

## Day-of checklist

- [ ] VPS healthy (`tailscale ping temporal-dev` from an instructor machine)
- [ ] Temporal UI loads
- [ ] Aperture healthy (test completion still works)
- [ ] Instruqt track published, invite links sent
- [ ] One test workflow runs end-to-end from an Instruqt VM
- [ ] Slides loaded

## Slides

The deck lives at `slides/slides.md` and runs under Slidev. See [Present the slides](slides.md) for installing pnpm, running the deck, presenter mode, and exporting to PDF.

## Backup plan

If the shared server fails mid-workshop:

1. Attendees run `temporal server start-dev` locally on their Instruqt VM.
2. They edit `~/.config/temporalio/temporal.toml` to set the `default` profile address to `localhost:7233` and run with `TEMPORAL_PROFILE=default`.
3. The Aperture integration still works. LLM calls route through Tailscale independently.
4. They lose the shared-UI "see everyone's workflows" experience, but every exercise still completes.

Rehearse this before the event. Switching an attendee over under pressure is frustrating if it's the first time you've tried.
