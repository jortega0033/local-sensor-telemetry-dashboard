# CLAUDE.md

Working instructions for this repository. Read this file first, then the documents it links to. It does not repeat them.

## What this is

An offline time-series workbench for one developer investigating synthetic machine telemetry: generate a run, ingest the readings into a local append-file engine, inspect a time window in Chrome as a chart plus an accessible table, restart, and see the same window again. Go standard library only, one executable, `127.0.0.1` only.

**Nothing is implemented yet.** The repository holds the plan (README, PRD, technical specification) and the backlog in epic #1. Every command below is planned and marked with the ticket that introduces it. Do not describe a planned thing as working, and do not claim a passing baseline: there is none.

## Sources of truth

| Topic | Document |
|---|---|
| Product scope, design, deterministic fixture, epic text | [docs/PRD.md](docs/PRD.md) |
| Technical design (same design sections as the PRD, plus fixtures) | [docs/TECHNICAL_SPEC.md](docs/TECHNICAL_SPEC.md) |
| Planned package map, freeze points, planned versus implemented | [docs/architecture.md](docs/architecture.md) |
| Backlog | Epic issue #1 and its child tickets on GitHub. The epic's table is authoritative for order, dependencies, parallel tracks and estimates; tickets state their own owner-decision defaults. |
| Loop setup and its limits | [LOOP.md](LOOP.md), [loop-constraints.md](loop-constraints.md) |
| Project subagents, origin and local changes | [.claude/agency-agents/README.md](.claude/agency-agents/README.md) |
| Ticket workflow and the step at which each role runs | [.claude/skills/implement-ticket/SKILL.md](.claude/skills/implement-ticket/SKILL.md), section "Agent roles" below |

The PRD and the technical specification repeat the same design sections. When they differ, the technical specification wins for design and the epic and tickets win for scope; a ticket that changes a contract fixes both documents in the same pull request.

## Stack and constraints

- Go module, **no external dependencies**: standard library only, in the application and in the tests. No package download in any launch, fixture or test step. The minimum Go version is an open decision with a stated default (ticket #2).
- One executable serves the API and the packaged HTML, CSS and JavaScript on `127.0.0.1`. The synthetic generator and the benchmark are separate subcommands of it.
- Offline: the browser fetches same-origin assets and data only. The demo is verified with external browser requests blocked. No CDN, font, image or script from another host.
- Supported host: one Linux host. One writer process per data directory, held by an OS-backed exclusive lease; a second writer fails at startup.
- Native browser assets: SVG for the chart, an HTML table for keyboard and assistive technology, a visible UTC label, the unit next to values.
- Data: integer gauge readings, one stable unit per series. Event time is integer UTC Unix milliseconds, 0 through year 2100; scaled integer values with magnitude at most 1,000,000,000; identifiers 1 to 40 characters from a portable set.
- Hard caps (technical specification): 500 readings and 256 KiB per POST; 31-day maximum range; 1,000 buckets per series; 1,000,000 scanned samples; 1,024 retained fingerprints per stream.
- No throughput or latency claim is made until the benchmark ticket records a measurement.
- A Go version change, a dependency, an external request or a non-loopback listener needs the owner's approval first, not a pull request.

## Commands (all planned)

Go is not installed on the owner's PC and nothing is installed globally. **Go checks run in CI only** (`actions/setup-go`, pinned by commit SHA). Run `go version` first: if it is missing, do not install it, and do not claim you ran any command below.

| Purpose | Command | Status |
|---|---|---|
| Build | `go build -o sensordash ./cmd/sensordash` | planned: introduced by ticket #2 |
| Format check (empty output is good) | `gofmt -l .` | planned: introduced by ticket #2 |
| Static check | `go vet ./...` | planned: introduced by ticket #2 |
| Tests | `go test ./...` | planned: introduced by ticket #2 |
| Run the server | `./sensordash serve --data-dir ./data --addr 127.0.0.1:8765` | planned: introduced by ticket #2 (flag names are proposals) |
| Print the synthetic fixture batch | `./sensordash gen --fixture` | planned: introduced by ticket #2 |
| Fault-injection tests (child process, controlled crashes) | `go test ./internal/store/... -run Crash` | planned: introduced by ticket #3 |
| Chrome browser tests (`CHROME_BIN` or `google-chrome` on `PATH`) | `go test ./internal/browsertest/...` | planned: introduced by ticket #10 |
| Benchmark with fixed seed and presets | `./sensordash bench --preset <name> --seed <n>` | planned: introduced by ticket #7 |

Tickets #8 to #10 are the new tickets of the epic rewrite (validation and request guards, single-writer lease and startup, dashboard). If they were created under other numbers, fix this table and the epic.

## Architecture boundaries

The planned package map is in [docs/architecture.md](docs/architecture.md). The technical specification names the boundaries, and they stay separate: **format validation, ingestion, storage, aggregation and HTTP**. In short:

- The format code owns frame encoding, checksums and versions and has no knowledge of HTTP.
- Ingestion owns validation, identity (sequence, fingerprints, high-water marks) and the commit, acknowledgment and duplicate policy.
- Storage owns the append file, the lease, startup verification, segments, the manifest and retention. Nothing else touches the data directory, and no path is ever derived from an HTTP value.
- Aggregation owns windows, buckets and merging across active and sealed data. It reads one published commit.
- HTTP owns routes, request guards (Host and Origin on mutations, capped JSON body), the error envelope and status codes. It never decodes the file format.
- The independent reference calculator used by tests must not import engine helpers.

Invariants every change must keep (technical specification, "Ingestion, acknowledgment, and recovery" and "Time and identity"): never acknowledge before `File.Sync` succeeds; a same-payload retry returns the original committed identity and adds no row; a changed payload for a retained fingerprint is 409 with no mutation; an uncertain write puts the writer in recovery-required state; readers never see a partial batch; decoders check limits before allocating and check overflow before adding deltas.

## Git workflow

- `main` is the default branch. Work on one branch per ticket: `<issue-number>-<short-slug>`. Never commit to `main`.
- Small commits with imperative subjects and a type prefix (`feat:`, `fix:`, `docs:`, `test:`, `chore:`).
- Order for every ticket: local evidence, Code Reviewer, Reality Checker, then a draft pull request that says `Part of epic #1` and `Closes #<ticket>`, then CI on that draft, then `gh pr ready` only when CI is green.
- Do not push or open a pull request unless the owner asked in the current session. Only the owner merges and releases. Do not push to `main`, force-push, or merge your own pull request.
- Never commit a data directory, `.env*`, or a built `sensordash` binary.
- Do not change `.claude/agents/`, `LOOP.md`, `loop-constraints.md` or `loop-budget.md` as part of a feature ticket.
- Open owner decisions on a ticket are settled with the ticket's stated default. Where a ticket states none, take the conservative option and say so in the pull request.

## Evidence before a ticket is called done

Paste real output into the pull request; never summarise a run you did not make.

1. For each acceptance criterion of the issue: the name of the test that proves it, or the exact step and what was observed.
2. Go results come from CI: the link to the CI run on the draft pull request and the result of `gofmt -l .`, `go vet ./...` and `go test ./...` in it, with the Go version. If Go is available locally, also paste the local output; if it is not, say so. Because there is no local Go run, the gates judge the diff, the tests written for each criterion and that CI runs them; a ticket is not done on a red or missing CI run.
3. For an HTTP change: one request and response pair per success and per error status, from a server on a scratch data directory (CI test output counts).
4. For a durability or recovery change: the fault point exercised, the oracle result (acknowledged rows survive, unacknowledged are absent or fully present, no duplicates), and the filesystem and fault model used.
5. For a UI change: the Chrome test run, plus the manual steps the ticket marks as owner checks; the page makes no external request.
6. `git status` shows no stray files, and no dependency, external URL or non-loopback listener was added.
7. Anything not verified, said plainly.

## Loop

Issue triage runs at L1: report-only, started by hand, no scheduler, no GitHub writes. See [LOOP.md](LOOP.md).

## Ticket implementation: agent roles run automatically

When asked to implement, continue or finish a GitHub issue, follow the playbook in [`.claude/skills/implement-ticket/SKILL.md`](.claude/skills/implement-ticket/SKILL.md). Claude Code loads it on its own for ticket work; `/implement-ticket <issue number>` starts it by hand.

Delegate to the project subagents in `.claude/agents/` without waiting to be asked:

- **Before coding:** Software Architect, as their descriptions say.
- **While coding:** the implementing role that matches the files being changed.
- **After coding:** each specialist reviewer whose trigger matches the change.
- **On every ticket, before it is called done:** Code Reviewer, then Reality Checker.
- The session that wrote a change does not count as its reviewer.
- List the roles that ran and their verdicts in the pull request description under "Agent roles run".

Each role's `description` ends with a "Use proactively ..." sentence, which is what makes Claude Code delegate to it without being told.
