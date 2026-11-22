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
| How a ticket goes from the queue to a ready pull request, who may do what | "Delivery workflow" below |
| UI and copy rules | [docs/ui-quality.md](docs/ui-quality.md) |
| Decisions the owner made, with dates | [docs/owner-decisions.md](docs/owner-decisions.md) |

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
| Chrome browser tests (binary from `CHROME_BIN`, else `google-chrome` on `PATH`; defined by ticket #10) | `go test ./internal/browsertest/...` | planned: introduced by ticket #10 |
| Benchmark with fixed seed and presets | `./sensordash bench --preset <name> --seed <n>` | planned: introduced by ticket #7 |

Tickets #8 to #10 are the new tickets of the epic rewrite (validation and request guards, single-writer lease and startup, dashboard). They are created in the order #8, #9, #10, before any pull request is opened in this repository, because issues and pull requests share one number sequence. If they got other numbers, fix every file that names them (`grep -rnE "#(8|9|10)\b" CLAUDE.md docs .claude/agents .claude/skills/next-ticket`) and the epic.

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
- Do not push or open a pull request unless the owner asked in the current session; starting `/next-ticket` is that request for one ticket. Only the owner merges and releases. Do not push to `main`, force-push, or merge your own pull request.
- Never commit a data directory, `.env*`, or a built `sensordash` binary.
- Do not change `.claude/agents/`, `LOOP.md`, `loop-constraints.md` or `loop-budget.md` as part of a feature ticket.
- Open owner decisions on a ticket are settled with the ticket's stated default. Where a ticket states none, take the conservative option and say so in the pull request.

## Evidence before a ticket is called done

Paste real output into the pull request; never summarise a run you did not make.

1. For each acceptance criterion of the issue: the name of the test that proves it, or the exact step and what was observed.
2. Go results come from CI, and green CI on the current diff is a condition for `gh pr ready`: the link to the CI run on the draft pull request (jobs `test` and `go`, and `bench` once it exists; Go steps run and not skipped) and the result of `gofmt -l .`, `go vet ./...` and `go test ./...` in it, with the Go version. If Go is available locally, also paste the local output; if it is not, say so. Because there is no local Go run, the gates judge the diff, the tests written for each criterion and that CI runs them; a ticket is not done on a red or missing CI run.
3. For an HTTP change: one request and response pair per success and per error status, from a server on a scratch data directory (CI test output counts).
4. For a durability or recovery change: the fault point exercised, the oracle result (acknowledged rows survive, unacknowledged are absent or fully present, no duplicates), and the filesystem and fault model used.
5. For a UI change: the Chrome test run, plus the manual steps the ticket marks as owner checks; the page makes no external request.
6. `git status` shows no stray files, and no dependency, external URL or non-loopback listener was added.
7. Anything not verified, said plainly.

## Agent roles

Project subagents are in `.claude/agents/`: ten Agency Agents roles and `loop-verifier`, which belongs to the Loop setup. Each Agency Agents role file starts with a section "In this repository", written for this project: the planned stack and hard rules that bind the role, the files it owns or reviews, its commands or the output it must be given, a checklist of this project's rules, the parts of the generic text that do not apply, and what it hands back. That section and this file override the generic text below it. Sources, licence, local changes and the name collision with a user-level install are in [.claude/agency-agents/README.md](.claude/agency-agents/README.md).

When asked to implement, continue or finish a GitHub issue, follow the playbook [`.claude/skills/implement-ticket/SKILL.md`](.claude/skills/implement-ticket/SKILL.md). Claude Code loads it on its own for ticket work; `/implement-ticket <issue number>` starts it by hand. Delegate to the roles without waiting to be asked: each role's `description` ends with a "Use proactively ..." sentence, which is what makes Claude Code delegate. The playbook gives the exact trigger of each role.

| When | Roles | Access |
|---|---|---|
| Before coding | Software Architect (module boundary, data model, file or manifest format, public API or dependency) | `Read, Grep, Glob` |
| While coding | Backend Architect (Go engine and API), Frontend Developer (`web/`) | can edit |
| While coding | Database Reliability Engineer (fsync order, recovery, manifest, retention, quota) | can edit |
| After coding, when the trigger matches | Application Security Engineer (untrusted input, Host and Origin checks, decoders), Accessibility Auditor (dashboard, table, keyboard), Performance Benchmarker (the benchmark ticket) | `Read, Grep, Glob` |
| After coding, when the trigger matches | API Tester (new or changed route, status code or contract; writes Go tests) | can edit |
| Gate on every ticket | Code Reviewer on the full diff, until no material finding is left | `Read, Grep, Glob` |
| Acceptance gate on every ticket, after Code Reviewer | Reality Checker on each acceptance criterion and the recorded evidence, and that CI is configured to run the checks (no CI result is needed yet) | `Read, Grep, Glob` |

- The two gates run on every ticket, whatever its size. A ticket is not called done before both have passed.
- Roles limited to `Read, Grep, Glob` cannot run commands. Give them the diff, the acceptance criteria and whatever command output exists. Go is not run locally, so say so instead of implying a run.
- The session or role that wrote a change does not count as its reviewer.
- If a role is missing or cannot be delegated to, the session does the before-coding, while-coding and specialist steps itself and says so. It cannot pass the two gates on its own change: do not push or open the pull request. Stop, report that the gates could not run, and leave the run for the owner. This overrides the playbook's sentence "do that step yourself" for the two gates.
- In the pull request description, list under "Agent roles run" each role that ran, what it was asked and its verdict, and any step done without a role.

## Delivery workflow

How a ticket reaches the owner. Each piece is described in its own file; this part only links them.

- **Start:** the owner runs `/next-ticket` ([`.claude/skills/next-ticket/SKILL.md`](.claude/skills/next-ticket/SKILL.md)); it picks one leaf ticket (never a `tracking` parent or a ticket with open sub-issues; the last sub-issue's pull request also closes the parent), creates the branch and worktree, and runs the playbook [`.claude/skills/implement-ticket/SKILL.md`](.claude/skills/implement-ticket/SKILL.md). Nothing is scheduled.
- **Gate hook:** [`.claude/hooks/role_gate.py`](.claude/hooks/role_gate.py) blocks `gh pr create` and `gh pr ready` until Code Reviewer and Reality Checker ran after the last code edit. `ROLE_GATE=off` in the environment Claude Code is started from switches it off for that session; only the owner does that.
- **Permissions:** [`.claude/settings.json`](.claude/settings.json). Agent configuration edits ask in the main checkout and in ticket worktrees; CI, `CLAUDE.md`, `AGENTS.md` and Loop-file edits ask only in the main checkout; in a ticket worktree they are limited by the playbook's scope rule and reviewed by the owner at merge (shell-made edits are not covered, so a ticket never changes agent configuration); force pushes, pushes to `main`, `gh pr merge`, releases and reads of `.env` are denied.
- **CI:** [`.github/workflows/ci.yml`](.github/workflows/ci.yml). Jobs `test` and `go` (and `bench` once ticket #7 adds it) must be green on the current diff before a pull request is marked ready; a failed `bench` job blocks `gh pr ready`, while a `bench` job that ends green with "no benchmark paths changed" does not. Until the first sub-issue of ticket #2 adds `go.mod`, the Go steps are skipped behind a guard that the same pull request removes, so no ticket pull request runs with the Go steps skipped. The agent does not poll: it ends its turn with "Draft pull request #N is waiting for CI" and reads `gh pr checks` once when woken. Reviewers read the diff and evidence from `.claude/state/review/`. An acceptance step that needs a person is listed under "Owner checks" in the pull request.
- **Pull request:** [`.github/pull_request_template.md`](.github/pull_request_template.md). Tickets use [`.github/ISSUE_TEMPLATE/implementation-ticket.md`](.github/ISSUE_TEMPLATE/implementation-ticket.md).
- **Labels:** `status:ready`, `status:blocked`, `status:in-progress`, `needs-decision`, `tracking` (parent container, never worked; its sub-issues are the work items).
- **Changelog:** every ticket pull request adds one line under Unreleased in [CHANGELOG.md](CHANGELOG.md).
- **UI and copy:** [docs/ui-quality.md](docs/ui-quality.md); the `ui-ux-pro-max` and `unslop` skills are in `.claude/skills/`.
- **Decisions:** [docs/owner-decisions.md](docs/owner-decisions.md). Security reports: [SECURITY.md](SECURITY.md).
- **Merge and release:** only the owner merges and releases ([docs/RELEASING.md](docs/RELEASING.md)). An agent stops at a ready pull request.

## Loop

Issue triage runs at L1: report-only, started by hand, no scheduler, no GitHub writes. See [LOOP.md](LOOP.md).
