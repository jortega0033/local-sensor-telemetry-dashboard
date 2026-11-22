# Releasing

Only the owner releases. Agents never tag, publish or create a GitHub release.

## What a release is here

A release is a Git tag on a commit of `main`, `vMAJOR.MINOR.PATCH`. Nothing is published and nothing is deployed: the product is one Go executable built from a clone (`go build -o sensordash ./cmd/sensordash`) and run on `127.0.0.1`. No release exists yet, and none should be tagged before the last ticket of epic #1 (the reproducible report and clean offline demo, #7) is merged.

## Checks before tagging

1. CI is green on the tip of `main` for both jobs (`test` and `go`) on Linux and Windows. On Linux every Go step runs and none is skipped (the guard on `go.mod` is gone after ticket #2); on Windows `go test` is skipped by design in `ci.yml` and the build and vet steps run.
2. `gofmt -l .` prints nothing, and `go vet ./...`, `go build ./...` and `go test ./...` pass in that CI run, with `GOPROXY=off`.
3. `go.mod` has no `require` line and there is no `go.sum`.
4. The offline demo steps in the README were run by the owner in Chrome against the built binary.
5. `docs/benchmark.md` holds numbers from a linked CI run (ticket #7).
6. `CHANGELOG.md`: move the `Unreleased` entries under a new version heading with the date.

## Tagging

1. In a pull request, update `CHANGELOG.md`; the owner merges it.
2. On the merge commit: `git tag -a vX.Y.Z -m "vX.Y.Z"` and `git push origin vX.Y.Z`.
3. Optional: create a GitHub release from the tag with the changelog section as its text.
