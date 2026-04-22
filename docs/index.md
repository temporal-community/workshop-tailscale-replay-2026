# Securing AI Applications with Tailscale and Temporal

A hands-on workshop that teaches three ideas by building one thing:

- **Durability** - Temporal orchestrates an agentic loop so the LLM can crash, retry, and resume without losing state.
- **Zero-config networking** - Tailscale connects every machine in the workshop over an encrypted mesh, no VPN concentrator, no firewall rules.
- **Shared-key safety** - Aperture acts as an API gateway so one OpenAI key serves the whole room with per-identity rate limiting.

## Who this site is for

This site is the **community guide**, a reference for anyone who wants to **run the workshop themselves**, whether that's a conference track, an internal training, or a weekend hack with friends.

If you are *attending* a session, the [GitHub README](https://github.com/temporal-community/workshop-tailscale-replay-2026) and the in-environment instructions have everything you need during the workshop itself. You don't need this site.

## Pick a path

<div class="grid cards" markdown>

- **I want to understand what we're building**

    ---

    Read the [workshop overview](workshop-overview.md) and the [architecture](architecture.md) to see the pieces and why they fit together.

- **I want to try it on my own machine**

    ---

    [Try it locally](try-it-locally.md) runs the full Tailscale + Temporal stack on your own machine, on your personal tailnet. Only Aperture is swapped (direct OpenAI key).

- **I'm teaching a cohort via Instruqt**

    ---

    [Running with Instruqt](run-with-instruqt.md) covers VM requirements, track configuration, and day-of checklists.

- **I'm teaching a cohort without Instruqt**

    ---

    [Running without Instruqt](run-without-instruqt.md) covers onboarding attendees onto your tailnet, Aperture alternatives, and the pre-event checklist.

- **I need to stand up the shared infrastructure**

    ---

    [Infrastructure](infrastructure.md) covers the persistent VPS that runs `temporal-ts-net` and the Tailscale and Aperture pieces behind it.

- **I just want to run the slides**

    ---

    [Present the slides](slides.md) covers installing pnpm, running the deck locally, presenter mode, and exporting to PDF.

</div>

## Source material

- **Workshop repo** - [temporal-community/workshop-tailscale-replay-2026](https://github.com/temporal-community/workshop-tailscale-replay-2026)
- **`temporal-ts-net`** - [temporal-community/temporal-ts-net](https://github.com/temporal-community/temporal-ts-net)
- **Temporal SDKs** - [Python](https://docs.temporal.io/develop/python), [Go](https://docs.temporal.io/develop/go)
- **Tailscale** - [docs](https://tailscale.com/kb). **Aperture** - [docs](https://docs.tailscale.com/aperture)
