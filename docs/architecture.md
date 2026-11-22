# Implementer map

Status: **planned, not implemented.** No Go file exists. This page maps the boundaries that [docs/TECHNICAL_SPEC.md](TECHNICAL_SPEC.md) already requires onto a package layout so that parallel tickets do not collide. It adds no requirement; the specification wins where they differ. Ticket #2 creates the layout and must correct this page to what it built, marking each row implemented.

## Planned package map

Module path: `github.com/jortega0033/local-sensor-telemetry-dashboard`. The names are proposals fixed by ticket #2.

| Boundary (spec, "Runtime and boundaries") | Package | Owns | Created by |
|---|---|---|---|
| Entry point | `cmd/sensordash` | Subcommands `serve`, `gen`, `bench`; flag parsing; startup message naming data root and bound address | #2 (`serve`, `gen`), #7 (`bench`) |
| Format validation | `internal/format` | Header, version, length-framed batches, CRC-32, decoder limits, later segment encoding | #2, extended by #3 and #5 |
| Ingestion | `internal/ingest` | Batch validation, identity rules, fingerprints, high-water marks, commit and duplicate policy | #2, extended by #8 and #3 |
| Storage | `internal/store` | Append file, `File.Sync`, lease, startup verification, recovery-required state, segments, manifest, retention, quota | #2, extended by #9, #3, #5 and #6 |
| Aggregation | `internal/query` | Half-open windows and buckets, count, sum, min, max, cross-block merge, scan and range caps | #2 (minimal), #4 |
| HTTP | `internal/httpapi` | Routes, Host and Origin checks, body cap, error envelope, status codes, static assets | #2, extended by #8 and #4 |
| Browser assets | `web/` (embedded with `embed`) | HTML, CSS, JavaScript, SVG chart, accessible table | #2, #10 |
| Synthetic data | `internal/gen` | Deterministic generator behind `gen` and `bench` | #2, #7 |
| Test oracle | `internal/refcalc` | Independent in-memory row filter and bucket calculator; imports no engine package | #4 |
| Browser tests | `internal/browsertest` | Headless Chrome runs against a real server | #10 |

Dependency direction: `httpapi` depends on `ingest` and `query`; both depend on `store`; `store` depends on `format`. `format` imports nothing from the module. `refcalc` imports nothing from the module.

## Freeze points

Parallel tickets rely on these shapes, which ticket #2 fixes and documents in this file:

- **Storage reader interface:** what `query` reads (one published commit, active log rows with series, producer, sequence, time, value).
- **HTTP result shapes:** the `POST /api/v1/readings` response (inserted and duplicate counts, new commit sequence, original commit sequence per retried row), the `GET /api/v1/query` response (per series, per bucket: count, sum, min, max, mean as sum and count; empty bucket has count and sum zero and null min, max and mean), and the error envelope `{"error": {"code", "message"}}`.
- **Append format version 1:** header and frame layout.

A ticket that changes a frozen shape updates this page and every dependent ticket in the same pull request.

## What is verified where

Go is verified in CI only (see [CLAUDE.md](../CLAUDE.md)). Crash-path tests live with `store`, browser fixtures with `browsertest`, and the fixture arithmetic from the specification is asserted against `refcalc` and the engine separately.
