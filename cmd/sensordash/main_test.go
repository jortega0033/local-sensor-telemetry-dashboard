package main

import (
	"bytes"
	"strings"
	"testing"
)

func runArgs(args ...string) (int, string, string) {
	var out, errb bytes.Buffer
	code := run(args, &out, &errb)
	return code, out.String(), errb.String()
}

func TestHelpAndVersion(t *testing.T) {
	code, out, _ := runArgs("--help")
	if code != 0 || !strings.Contains(out, "sensordash serve --data-dir") {
		t.Fatalf("--help: code %d, output %q", code, out)
	}
	code, out, _ = runArgs("--version")
	if code != 0 || !strings.Contains(out, version) {
		t.Fatalf("--version: code %d, output %q", code, out)
	}
}

func TestUnbuiltSubcommandsFailClearly(t *testing.T) {
	for name, ticket := range pending {
		code, _, errOut := runArgs(name)
		if code == 0 || !strings.Contains(errOut, "not implemented yet: see "+ticket) {
			t.Errorf("%s: code %d, stderr %q", name, code, errOut)
		}
	}
}

func TestUnknownCommand(t *testing.T) {
	if code, _, _ := runArgs("nope"); code == 0 {
		t.Fatal("unknown command must exit non-zero")
	}
}
