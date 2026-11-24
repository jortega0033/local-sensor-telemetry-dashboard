# Local Sensor Telemetry Dashboard

A planned offline project. Implementation has not started. The repository currently contains the product plan, technical specification and implementation backlog.

## Program goal

Build a full time-series database and monitoring system in Go's standard library: compressed columnar storage, a query language, anomaly detection, alerting and replication, able to hold billions of readings.

The program is planned ([#32](https://github.com/jortega0033/local-sensor-telemetry-dashboard/issues/32)). Phase 1 is the current epic [#1](https://github.com/jortega0033/local-sensor-telemetry-dashboard/issues/1) and the only phase with tickets.

**Score.** Readings per second and bytes per reading on a generated industrial workload, with every answer equal to a naive oracle.

| Phase | Content | Hours (rough) | Issue |
|---|---|---|---|
| 1 | The current epic and its tickets | 63 to 97 | [#1](https://github.com/jortega0033/local-sensor-telemetry-dashboard/issues/1) |
| 2 | Storage | 140 to 240 | [#33](https://github.com/jortega0033/local-sensor-telemetry-dashboard/issues/33) |
| 3 | Query engine | 130 to 220 | [#34](https://github.com/jortega0033/local-sensor-telemetry-dashboard/issues/34) |
| 4 | Understanding | 120 to 210 | [#35](https://github.com/jortega0033/local-sensor-telemetry-dashboard/issues/35) |
| 5 | Cluster | 170 to 290 | [#36](https://github.com/jortega0033/local-sensor-telemetry-dashboard/issues/36) |
| 6 | Protocols and view | 120 to 200 | [#37](https://github.com/jortega0033/local-sensor-telemetry-dashboard/issues/37) |

About 740 to 1260 engineer hours in total: rough estimates, not measurements. The numbers in the score and the finish lines are targets until a run shows them. Only phase 1 has tickets to work; the sub-issues of a later phase are written when it starts. No single session is expected to reach the end.

## Plan

- [Full product plan and acceptance criteria](docs/PRD.md)
- [Technical specification and verification fixtures](docs/TECHNICAL_SPEC.md)
- [Implementation epic](https://github.com/jortega0033/local-sensor-telemetry-dashboard/issues/1)

## Current status

The fixtures describe expected results for future checks. No working baseline or passing implementation tests are claimed. Required runtime versions and offline dependencies must be checked before the first build.

## Licensing

A license for this project has not been selected. Any upstream reference or adapted code must retain its applicable notices; see the plan.
