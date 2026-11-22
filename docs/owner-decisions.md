# Owner decisions

A dated record of the decisions the owner made for this repository, so the rules do not depend on anyone's memory of a conversation. Changing any entry needs the owner. An agent that finds a decision blocking a ticket stops and says so.

All entries are dated 2026-11-22.

| # | Decision | Date | Implemented in |
|---|---|---|---|
| 1 | The instructions, tooling and playbook live on the default branch (`main`). | 2026-11-22 | `CLAUDE.md`, `.claude/`, `.github/` |
| 2 | Persona agents run automatically during ticket work and carry instructions tailored to this project. | 2026-11-22 | `.claude/agents/` ("Use proactively" sentences, "In this repository" sections), `.claude/agency-agents/README.md`, `CLAUDE.md` "Agent roles" |
| 3 | Every open ticket is about app implementation; tooling is installed, not ticketed. | 2026-11-22 | `.claude/`, `.github/`, `docs/`, `CLAUDE.md` "Delivery workflow" |
| 4 | Agents may take a ticket up to a ready pull request. Only the owner merges and releases. | 2026-11-22 | `CLAUDE.md` "Git workflow" and "Delivery workflow", `.claude/skills/next-ticket/SKILL.md`, `.claude/skills/implement-ticket/SKILL.md`, `docs/RELEASING.md`, `.claude/settings.json` (`gh pr merge` and `gh release create` denied) |
| 5 | No scheduler. Ticket work starts only on the owner's go signal (`/next-ticket`), because the repository is used to train AI models. | 2026-11-22 | `.claude/skills/next-ticket/SKILL.md` (`disable-model-invocation: true`), `LOOP.md` |
| 6 | A review gate hook blocks opening or readying a pull request until the gate roles have run after the last code edit. | 2026-11-22 | `.claude/hooks/role_gate.py`, `.claude/settings.json` (`hooks`), `.gitignore` (`.claude/state/`) |
| 7 | Loop Engineering stays at L1, report-only and unscheduled. | 2026-11-22 | `LOOP.md`, `loop-constraints.md`, `loop-budget.md` |
| 8 | Continuous integration runs on Linux and Windows; the default branch requires a pull request and green checks. | 2026-11-22 | `.github/workflows/ci.yml`; branch protection is a repository setting applied by the owner's orchestrator, not a file |
| 9 | Permission rules: documented commands are allowed; force pushes, pushes to the default branch, merges, releases and reads of secrets are denied; edits to CI, agent configuration and Loop files ask first. | 2026-11-22 | `.claude/settings.json` (`permissions`) |
| 10 | Open decisions on tickets are settled with the ticket's stated default; where a ticket states none, the conservative option is taken and flagged in the issue. | 2026-11-22 | The list below; the issues themselves |
| 11 | Go is not installed on the owner's PC and nothing is installed globally. Go is verified in CI only (`actions/setup-go` pinned by commit SHA), like Rust in the owner's other repositories. The Go module has no external dependencies, the supported host is Linux, and CI is the only place where the Go version from `go.mod` is installed. | 2026-11-22 | `.github/workflows/ci.yml` (job `go`, guarded until ticket #2 adds `go.mod`), `CLAUDE.md` "Commands (all planned)", `AGENTS.md` |
| 12 | The ten UI rules in `docs/ui-quality.md` are binding, and the `ui-ux-pro-max` and `unslop` skills are installed in the repository and used for UI work and for copy. | 2026-11-22 | `docs/ui-quality.md`, `.claude/skills/ui-ux-pro-max/`, `.claude/skills/unslop/`, `.claude/skills/implement-ticket/SKILL.md`, acceptance criteria of #2, #6 and #10 |

## Open ticket decisions settled by default (entry 10)

Each is settled with the default below unless the owner changes it in the issue. None blocks a ticket.

| Issue | Decision | Default | Status |
|---|---|---|---|
| #2 | Go minimum version | `go 1.26` in `go.mod` (the older of the two supported releases on 2026-11-22; Go 1.27.4 and 1.26.11 are current). CI installs the version from `go.mod`. | settled by default |
| #2 | Module path and package layout | `github.com/jortega0033/local-sensor-telemetry-dashboard`, layout in `docs/architecture.md`. | settled by default |
| #2, #9 | Platform | Linux is the supported host. Non-Linux builds compile with a lease stub that returns "unsupported platform". | settled by default |
| #3 | Data-directory filesystem and crash tests | CI scratch directory on the runner's local disk; the fault model is a controlled process crash at named points, not power loss; the report prints the filesystem type. | settled by default |
| #8 | Identifier set, caps and bounds | Identifiers `[A-Za-z0-9._-]`, 1 to 40 characters; at most 100 series and 100 producers; event time 0 to 4102444800000 inclusive (2100-01-01T00:00:00Z); value magnitude at most 1,000,000,000; sequence 1 to 9007199254740991. | settled by default |
| #2, #3, #4, #6, #8, #9 | Status codes the plan leaves open | `POST /api/v1/readings` success is 200, also for an all-duplicate retry; 413 for more than 500 readings or more than 256 KiB; 400 for every other invalid request including query-cap rejections and capacity violations; 409 for sequence, fingerprint, time-regression and below-cutoff conflicts; 503 recovery-required and transient storage failure; 507 quota. Precedence for one row, first match wins: exact retry with a retained fingerprint is a duplicate (original commit identity, no new row, reading stays invisible if expired); retained fingerprint with other content is 409 `sequence_conflict`; unseen sequence at or below the high-water mark or outside the 1,024 fingerprint window is 409 `stale_sequence`; earlier event time is 409 `time_regression`; a new sequence with a reading older than the cutoff is 409 `below_cutoff`. A batch with any rejected row writes nothing. | settled by default |
| #2, #10 | Mean rounding | The mean is derived from sum and count; shown in units with three decimals, rounded half away from zero, in integer arithmetic. | settled by default |
| #5 | Segment size and compression | Seal at 4,096 rows or 1 MiB of active log, whichever comes first; varint scheme as in the technical specification; not retuned in this epic. | settled by default |
| #10 | Browser test method | Headless Chrome driven over `--remote-debugging-pipe` from Go tests, Linux CI only; skipped when Chrome is absent unless `CI=true`. | settled by default |
| #6 | Retention and quota | Both off by default (`--retention`, `--max-data-bytes`); boundary-segment rewrite is not part of this epic. | settled by default |
| #7 | Benchmark | Runs on the GitHub-hosted `ubuntu-latest` runner in a `bench` job of `ci.yml` that runs its benchmark steps on pull requests touching the benchmark (a changed-files step using `git diff --name-only` with full history decides, because GitHub has no job-level path filter and `on.pull_request.paths` would also filter the required jobs; with no benchmark path changed the job ends green) and on `workflow_dispatch` for later reruns (a manual workflow cannot run before it is on the default branch, so the numbers in `docs/benchmark.md` come from the pull request's own run on its head commit); presets `small`, `medium`, `large`; no sustained throughput goal is chosen. | settled by default |
| #2, #10 | Owner checks | Linux is the only supported host. The owner's PC is Windows without Go: owner checks run in WSL or on a Linux host, using the `sensordash-linux-amd64` binary that CI uploads as an artifact (ticket #2), so no Go is needed. | settled by default |
| none | Licence | Not selected. No ticket depends on it and none adds a LICENSE. | settled by default |
