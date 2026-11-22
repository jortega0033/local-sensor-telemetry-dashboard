# Local Sensor Telemetry Dashboard

## Product requirements

A developer investigating synthetic machine telemetry can generate a run, ingest its readings into a local time-series engine, and inspect a time range in Chrome. The chart and its adjacent table show when values changed, how many samples contributed, and which intervals were empty. The developer can restart the service and reopen the same window.

The first milestone is a working path from one generated stream through a durable append file and a localhost API to a browser table. It supports a fixed-width aggregation and an idempotent replay of the same batch. Later milestones add fault recovery, multiple streams, charting, compressed segments, retention, and benchmarking. The batch uses generated data, the Go standard library, and native browser assets. It runs on a single Linux host without external services.

The epic is complete when a clean checkout runs the synthetic ingest-to-dashboard demo, restart and fault fixtures preserve the specified consistency contract, compressed and active readings give identical query results, and an offline benchmark records actual resource use. Browser users can inspect the same values through the chart or an accessible table and distinguish empty, pending, and failed query states.

The scope covers integer gauge readings with one stable unit per series and one writer process per data directory. Device protocols, live device connections, arbitrary write ordering, multi-node replication, alerts, user accounts, and real telemetry belong to later work.

## Technical design

### Runtime and boundaries

Use a Go module with no external dependencies. Check the available version before setting the module's minimum. A single executable serves the API and packaged HTML, CSS, and JavaScript on 127.0.0.1. Put the synthetic generator and benchmark behind separate subcommands. The engine keeps format validation, ingestion, storage, aggregation, and HTTP concerns separate. Accept an explicit app-owned data directory, and hold an OS-backed exclusive lease for its lifetime. Fail startup if a second writer is already running. Release resources through the owning process rather than trusting a stale marker file.

The browser fetches same-origin assets and data only. Validate Host and Origin on mutations, take JSON with a capped body, and never derive a filesystem path from an HTTP value. Use native SVG for a semantic, scalable chart and an HTML table for keyboard and assistive technology access. Use a visible UTC label and expose the unit alongside the values. Ship an empty state and an opt-in fixture button. Name the chosen data root and bound address on startup. Verify the demo with external browser requests blocked.

### Time and identity

Represent event time as integer UTC Unix milliseconds and render ISO 8601 at the UI boundary. Limit input to nonnegative times through 2100 and scaled integer values with magnitude at most 1,000,000,000. Metric definitions declare the unit and scale. Values stay integers internally. Each sample carries series ID, producer ID, monotonically increasing sequence within that pair, event time, and value. Restrict identifiers to 1–40 characters and an explicit portable character set. Cap series and producer cardinality in configuration and validate capacity before admitting an entire batch.

An exact retry of an accepted sequence produces its original committed identity and adds no reading. If a retained fingerprint disagrees, return 409 without mutation. Reject unseen sequences below the stream high-water mark and new sequences with earlier event times. Sequence gaps are allowed and need no placeholder rows. Retain only the most recent 1,024 committed fingerprints per stream. Carry forward the high-water mark and last event time after data retention deletes old samples. An old retry outside that fingerprint window returns a stale-sequence conflict. This bounds index memory and prevents a retry from resurrecting expired data. Persist the window and high-water marks in a checkpoint when old append data is compacted. The retry horizon is count-based, so the UI must never promise an unlimited retry window.

### Ingestion, acknowledgment, and recovery

POST /api/v1/readings accepts at most 500 readings and 256 KiB of UTF-8 JSON. Reject unsupported fields, invalid types, unknown series definitions, and capacity violations before writing. The first version serializes the write path. Build tentative stream states for the whole batch, distinguish duplicates from new rows, and commit only when every row has passed validation. Return inserted and duplicate counts, a new commit sequence when readings were added, and each row's original commit sequence for retries. Use 400 for invalid input, 409 for sequence conflicts, 413 for oversize input, and retryable 503 for transient storage failure. Return 507 for a configured disk quota refusal. Browser reads use one published commit and never observe a partially accepted batch.

Store a versioned header followed by length-framed batches. A frame contains payload length, canonical payload, and CRC-32 of the payload. Reject unsupported versions, oversized frames, and corrupt committed data. Serialize append, validate complete writes, call File.Sync and wait for success, then publish the durable commit to readers and acknowledge the caller. A crash after persistence and before the response can leave the batch present. A same-payload retry must return the matching committed result without adding rows. Any short write or Sync failure with an uncertain result puts the writer in recovery-required state and blocks new writes. Restart or inspect the log before retrying writes. File.Sync is a file durability primitive, and the verification gate targets documented Linux filesystems and controlled process crashes. [Go File.Sync](https://pkg.go.dev/os#File.Sync) · [CRC-32](https://pkg.go.dev/hash/crc32)

Startup verifies lengths, checksums, stream sequence, and commit order. A provably incomplete frame at EOF can be truncated to the last verified offset, with a diagnostic. A checksum mismatch in a structurally complete frame, an impossible length, invalid version, or interior corruption stops writable startup. Preserve the bytes for investigation. Never acknowledge a record before durability. The test oracle records successful acknowledgments separately from the engine data directory and checks that every acknowledged batch survives restart, that unacknowledged attempts are absent or fully present, and that no batch appears twice. Distinguish process-crash evidence from untested hardware behavior.

### Query and display

GET /api/v1/query accepts one or more series IDs, integer start and end timestamps, and a whole-millisecond positive bucket width. Use half-open ranges and buckets. Anchor bucket boundaries at the requested start and clip the final bucket at the requested end. Reject end less than or equal to start, ranges beyond 31 days, requests above 1,000 buckets per series, or queries which would scan more than 1,000,000 samples. The engine enforces those caps before an unbounded operation. Window aggregation produces count, sum, minimum and maximum. An empty bucket has count and sum zero and null min, max, and mean.

Use exact signed 64-bit additions under the sample and value caps. Serialize mean as derived sum and count, with the UI showing three decimal places using a documented rounding rule. If multiple source blocks contribute to a bucket, merge their sums and counts before deriving a mean. Sort by series and then bucket start. Dashboard filters control series, UTC range, and bucket size. The chart skips empty points and does not connect a line across an empty interval. A newer request supersedes stale responses, and a failed refresh retains the last successful window with its label and an error.

### Segment compression and retention

Later, seal immutable segments in commit order. Within each segment group by series and keep timestamp deltas as unsigned varints, signed scaled values or signed value deltas as varints, and sequences as explicit monotonic deltas per stream. Segment metadata carries version, first and last time per series, row count, covered commit range, and checksum. Decoder limits precede allocation. Check bounds and overflow before adding deltas. The scan planner skips segments with no overlap, and a query snapshot sees a commit in either the active log or its sealed segment exactly once. [Go binary encodings](https://pkg.go.dev/encoding/binary)

Seal in a temporary file on the same filesystem, validate it, sync it, and publish via versioned manifest replacement. Sync the relevant directory when the target storage contract requires it. A manifest names complete segments, replay start, stream checkpoints, and retention watermark. If the process stops before the new manifest becomes current, replay from the old authoritative generation. Keep old referenced files until publication is durable and older reader snapshots no longer need them. Never infer current visibility from stray files. Run one last retention pass only after the replacement manifest has passed verification.

Retention uses an injected clock and a monotonic cutoff. Query visibility immediately excludes samples earlier than cutoff, including samples in a boundary segment. Non-overlapping old segments can be deleted after publication, and a boundary segment can be rewritten in a later no-data-loss compaction. Reads at cutoff remain visible. Reject batches with new readings older than cutoff, preserve recent retry fingerprints and the stream high-water marks, and expose the oldest queryable timestamp. A disk quota or slow flush changes ingestion status visibly without dropping acknowledged records.

## Deterministic fixture and verification

Use an original dataset anchored at 2026-01-01T00:00:00Z. Each series uses producer P and integer milliunits. The numbers in the table are offsets in milliseconds.

| Series | Sequence | Offset | Value |
|---|---:|---:|---:|
| A | 1 | 0 | 1200 |
| A | 2 | 59000 | 1800 |
| A | 3 | 60000 | -600 |
| A | 4 | 120000 | 2400 |
| B | 1 | 15000 | 500 |
| B | 2 | 75000 | -500 |

For a query over [0, 120000) with 60000 ms buckets, A has count 2, sum 3000, min 1200, max 1800, and mean 1500 in bucket 0. Its second bucket has count 1 and sum, min, max, and mean all equal to −600. Its whole-range sum is 2400 from three readings, so the mean is 800. The A reading at 120000 is excluded. B has sum 500 then −500, each with count 1. Extend the query to 180000 and B's third bucket is empty. Retry A sequence 2 with the same value and get one duplicate and no new row. Change its payload and get 409. With a cutoff of 60000, A sequence 2 is invisible and sequence 3 is visible.

Use an independent in-memory row filter and bucket calculator for fixture verification. The expected values above are derived from the input arithmetic. Validate them against the implementation later. Add generated edge fixtures for full commit plus lost response, truncated tail, interior corruption, retention boundary, negative values, old retries outside the fingerprint window, and identical values in cold and hot storage. Verify unknown versions fail before opening an ingest endpoint. For browser acceptance, verify the exact window and accessible table cold, after reload, after a service restart, and with external requests blocked.

## Epic E73 · Inspect generated sensor readings after a restart

Deliver an offline time-series workbench that durably ingests generated sensor streams, renders a chosen window, and demonstrates how aggregation and retention behave across restart and segment sealing. Acceptance requires passing the fixture, retry, crash-recovery, segment parity, and browser tests above, together with a reproducible benchmark report.

| Ticket | Scope | Dependencies | Estimate |
|---|---|---|---:|
| E73-1 | Run one durable ingest-to-table demo | — | 3–5 h |
| E73-2 | Harden retry, validation, and crash recovery | E73-1 | 3–5 h |
| E73-3 | Build exact windows and a browser dashboard | E73-1 | 3–5 h |
| E73-4 | Seal and read compressed segments | E73-2, E73-3 | 3–5 h |
| E73-5 | Add retention, limits, and safe reclamation | E73-4 | 2–4 h |
| E73-6 | Release gates and benchmark report | E73-3, E73-5 | 2–4 h |

**Remaining effort.** 16–28 engineering hours, including tests and integration. E73-2 and E73-3 can run in parallel after E73-1 freezes the storage-reader and HTTP result shapes. Assign crash-path tests to the storage lane and browser fixtures to the UI lane. The lead integrates both in E73-4 and runs the end-to-end verification in E73-6.

### E73-1 · Durable ingest-to-table demo

Implement the app start command, one synthetic producer, one series, a versioned append format, and a localhost flow from fixture submission to browser table. Use a fixed window initially and expose the seed fixture action only with explicit click. Verify A's first two buckets against the table, restart the same data directory, and verify the same browser output. Retry an unchanged batch twice and check the stored row count. The launch and fixture steps must run from a clean checkout without package downloads.

### E73-2 · Retry and crash recovery

Validate complete batches before mutation and enforce the commit, acknowledgment, and duplicate policies. Add deterministic child-process fault checkpoints before append, after partial bytes, after durable commit, and before acknowledgment. Verify observed acknowledged rows always reappear, partial tail is ignored with a diagnostic, interior corruption fails safely, and same-payload retries cannot duplicate a row. Check changed and stale sequence behavior and bounded fingerprint recovery. Start with a disk-backed scratch root and report the filesystem and fault model.

### E73-3 · Exact windows and dashboard

Implement the query window API, cross-block aggregation contract, and visible chart/table with series filters and empty buckets. Verify A and B together using the expected values and UTC labels. Add an empty interval and a complete refresh after an injected API error. Exercise range and bucket cap rejection and one race between two browser query requests. Recover the last good view without displaying it under a new failed request's label.

### E73-4 · Compressed segments and restart

Implement validated segment encoding and mixed active/sealed scans. Check malformed lengths, maliciously large declarations, overflow in decoded deltas, checksum mismatch, and non-sequential segment manifests. Compare active-only, sealed-only, and mixed results byte-for-byte at the normalized aggregate layer for the original fixture and a seeded 10,000-row stream. Simulate interruption before manifest publication and after new manifest publication. Verify complete rows and their retry metadata survive both points and never double-count.

### E73-5 · Retention and resource controls

Add controlled UTC time, visibility cutoff, manifest-backed cleanup, and quota error handling. Verify the exact 60000 ms boundary fixture, recover after interruption and server restart, and retain enough identity state to reject stale retries. Hold one live reader while publishing a replacement generation and verify the old files remain available for that reader until it finishes. Check that quota refusal leaves all prior acknowledgments intact, changes no query results, and is clearly shown in the dashboard.

### E73-6 · Reproducible report and clean offline demo

Create a deterministic seed and workload presets with stated series cardinality, batch size, durations, and query-window widths. Run three timed repeats on a stated host after a warm-up. Report accepted unique samples per second as unique committed readings divided by measured monotonic elapsed seconds. Specify whether HTTP and durability waits are included. Report commit and query p95 in milliseconds using nearest rank with sample counts, peak resident MiB, and stored data bytes per retained reading. Define stored data as the current WAL, segments, and active metadata, report logical bytes and optionally allocated disk bytes separately, and disclose any retained obsolete files. Compare uncompressed and compressed formats over the same dataset and settings. Confirm full fixture, fault, and Chrome runs on the same revision. Record results and the bottleneck before choosing a sustained throughput goal. The official inspiration title is a longer-term direction, with no measured throughput claimed in this plan.

## Open decisions before implementation

Confirm the supported Go version, data-directory filesystem and local crash-test permissions. Benchmark the chosen segment size and compression strategy once the uncompressed path exists. Keep the reference calculator independent of engine helpers. Start with uniform gauges and fixed integer quantities before adding more formats.

[Official project inspiration](https://docs.google.com/document/d/1R5MyfVvn033KoXZXgXP65qRAXnhcjtVvO_VF7IxUZaA/edit?tab=t.ef9a30wzb1z5)
