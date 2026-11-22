# Local Sensor Telemetry Dashboard: technical specification

Planning document. Implementation and verification remain future work.

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
