// Package store owns the append file, File.Sync, the exclusive lease, startup verification, recovery-required state, segments, the manifest, retention and quota. Nothing else touches the data directory.
package store
