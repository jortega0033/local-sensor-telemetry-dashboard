# Agency Agents roles used in this repository

The Markdown files in `.claude/agents/` whose names start with `engineering-`, `security-` or `testing-` are copies of role definitions from Agency Agents. They are project-local Claude Code subagents. The file `loop-verifier.md` in the same folder is not from Agency Agents; it belongs to the Loop Engineering setup described in [`LOOP.md`](../../LOOP.md).

## Source and licence

- Upstream: https://github.com/msitarzewski/agency-agents
- Pinned revision: `8ae6164602c95ca78828b148b9052460ee426885` (upstream `main`, committed 2026-11-09)
- Licence: MIT, "Copyright (c) 2025 AgentLand Contributors". The upstream licence text at the pinned revision is kept unchanged in [`LICENSE`](LICENSE) next to this file (SHA-256 `9a45258434d5cedf0af73c9ad4771373701225038d246c49219026c33677f66f`). It covers only the copied role files, not the rest of this repository.
- Fetched with `gh api -H "Accept: application/vnd.github.raw" "repos/msitarzewski/agency-agents/contents/<path>?ref=<revision>"`.

## Installed roles

"Upstream SHA-256" is the hash of the file as fetched at the pinned revision. Reviewer roles have one added frontmatter line, `tools: Read, Grep, Glob`, so they cannot edit files or run commands. Nothing else in any file was changed at installation.

| File in `.claude/agents/` | `name:` | Upstream path | Upstream SHA-256 | Local change |
|---|---|---|---|---|
| `engineering-code-reviewer.md` | Code Reviewer | `engineering/engineering-code-reviewer.md` | `7509fcc3ea1dda46511b2801996305bf3d4a125576f9655ff0b63df25602f446` | added `tools: Read, Grep, Glob` |
| `testing-reality-checker.md` | Reality Checker | `testing/testing-reality-checker.md` | `6d32fcdb114233e13902ec6372d50293b120e85d490b5e81d372c29808f988a1` | added `tools: Read, Grep, Glob` |
| `engineering-software-architect.md` | Software Architect | `engineering/engineering-software-architect.md` | `b85121691776147bacde26673e68fbbacd6a7f0816d1ca8e203c264d25ce1d30` | added `tools: Read, Grep, Glob` |
| `security-appsec-engineer.md` | Application Security Engineer | `security/security-appsec-engineer.md` | `a19f0bda46f4adca5737365dbc3c7ed93e9d93b71349536d92ff055dcf4e3b83` | added `tools: Read, Grep, Glob` |
| `testing-accessibility-auditor.md` | Accessibility Auditor | `testing/testing-accessibility-auditor.md` | `362935907af23839775b9c6162f38c0c5a740542ab4350cb5c5b5a0d9f9a4268` | added `tools: Read, Grep, Glob` |
| `testing-performance-benchmarker.md` | Performance Benchmarker | `testing/testing-performance-benchmarker.md` | `16a841f67392d084b773c6d9f05341d093e9017f3e116f6df0de59a5a3f59c34` | added `tools: Read, Grep, Glob` |
| `engineering-backend-architect.md` | Backend Architect | `engineering/engineering-backend-architect.md` | `ac717bde4848d7663fff8a3c7edc711fd119256de6a7ce38c036619372f42c99` | none |
| `engineering-database-reliability-engineer.md` | Database Reliability Engineer | `engineering/engineering-database-reliability-engineer.md` | `0096dc3c1b31e33814c509fda43853cc77cef6252bb94da527b833be13b53a5a` | none |
| `engineering-frontend-developer.md` | Frontend Developer | `engineering/engineering-frontend-developer.md` | `1ae4cc10b2e0163bd3ed921ad54080b865cd02e591becd8919b64b73fab06dad` | none |
| `testing-api-tester.md` | API Tester | `testing/testing-api-tester.md` | `ab91de9188d48bd9c417282c33673a37b066215e90e08c0f1b270d036084e3bb` | none |

## Why these roles

| Role | Use in this repository |
|---|---|
| Code Reviewer (read-only) | Review a diff for correctness, missed edge cases and broken rules from `CLAUDE.md`: acknowledgment before sync, duplicate handling, caps, overflow, decoder limits before allocation. |
| Reality Checker (read-only) | Compare claimed evidence with what was actually produced before a ticket or the epic is called done. Go checks run in CI only, so it also checks that a test exists for each acceptance criterion and that CI runs it. |
| Software Architect (read-only) | Check a ticket against the five boundaries (format, ingestion, storage, aggregation, HTTP) and the freeze points in `docs/architecture.md` before a data model, file format or interface changes. |
| Application Security Engineer (read-only) | Review where untrusted input meets the engine: Host and Origin checks, capped JSON, identifiers, decoder bounds, no path from HTTP input. |
| Accessibility Auditor (read-only) | Review the table, chart alternative, focus, labels and live status regions of the dashboard. |
| Performance Benchmarker (read-only) | Review the benchmark method and report of the final ticket: workload presets, repeats, nearest-rank p95, what the numbers include. |
| Backend Architect | Implement Go code for ingestion, storage, aggregation and the HTTP API. |
| Database Reliability Engineer | Implement and review crash recovery, fsync ordering, the manifest and segment publication, retention and quota behaviour. |
| Frontend Developer | Implement `web/`: native HTML, CSS, JavaScript and SVG. |
| API Tester | Write HTTP-level Go tests for the routes, status codes and error envelope. |

Most of these roles are written for larger stacks (frameworks, cloud services, PostgreSQL, Playwright, scanners). In this repository `CLAUDE.md` wins: Go standard library only, no package downloads, offline, `127.0.0.1` only.

Roles that were considered and left out: Codebase Onboarding Engineer (no code exists yet), Test Automation Engineer and Evidence Collector (their methods need Playwright or other downloads; browser tests use headless Chrome from the Go test suite), UI Designer and UX Architect (one small dashboard screen; the Frontend Developer and Accessibility Auditor cover it), Technical Writer and Developer Tooling Engineer (tooling is installed by the owner's delivery setup, not by tickets).

## Name collision with user-level agents

A developer who has installed Agency Agents in `~/.claude/agents/` already has agents with the same `name:` values. Claude Code gives the project-level definition precedence, so inside this repository the copies here are used, including the read-only tool limits. The upstream names were kept on purpose so that the roles stay recognisable and can be compared with upstream.

A Claude Code session that was already running when these files were added keeps the definitions it loaded at start. Start a new session to pick up the project copies.

## Check that the roles load

From the repository root, in a new session:

```text
claude -p "List the project-level subagents defined in .claude/agents with the tools each one may use. Then run the Reality Checker subagent and ask it which tools it has."
```

The six reviewer roles and `loop-verifier` must report exactly `Read, Grep, Glob`.

## Updating

Fetch the same paths at a newer revision, re-apply the `tools:` line to the six reviewer roles, and update the revision and hashes in this file in the same commit.

## Local change: automatic delegation (2026-11-22)

One sentence beginning "Use proactively" was appended to the `description:` line of every Agency Agents role in `.claude/agents/` (not to `loop-verifier`). Claude Code uses that field to decide when to delegate, so the roles now run during ticket work without being named.

A section headed "In this repository" was then inserted directly after the frontmatter of every Agency Agents role (not into `loop-verifier`), before the upstream text. It is project-specific text written for this repository, not upstream text: the planned stack and hard rules that bind the role, the files it owns or reviews, its commands or the output it must be given, a checklist of this project's rules with their sources, the parts of the generic text that do not apply here, and what it hands back. It ends with the sentence "This section and CLAUDE.md override the generic text below wherever they differ." The upstream text below it and the rest of the frontmatter are unchanged. The repository has no code yet, so the sections describe the planned stack, packages and commands and say so; they are corrected when ticket #2 builds the first version.

This supersedes any statement above that a role file is unmodified or that the `tools:` line is the only local change. To compare a role with the pinned upstream file, first remove the "Use proactively" sentence (in Backend Architect, Frontend Developer, API Tester, Performance Benchmarker and Reality Checker also the full stop that was added in front of it, because the upstream description has none), the section "In this repository" together with the blank line after it, and the `tools:` line where one was added. Checked on 2026-11-22: with these three removed, each of the ten files has the upstream SHA-256 listed in the table above. After fetching a newer upstream revision (see "Updating"), put all three back and check the section against the code as it then is.

The steps at which each role runs are in `.claude/skills/implement-ticket/SKILL.md`.
