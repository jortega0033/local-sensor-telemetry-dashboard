// Command sensordash is the single executable: it serves the API and the
// packaged page on 127.0.0.1, and hosts the synthetic generator and the
// benchmark as separate subcommands.
package main

import (
	"fmt"
	"io"
	"os"
)

const version = "0.0.0-scaffold"

const usage = `Usage:
  sensordash serve --data-dir <dir> --addr 127.0.0.1:<port>
  sensordash gen --fixture
  sensordash bench --preset <name> --seed <n>
  sensordash --version
  sensordash --help
`

// pending maps each subcommand to the ticket that implements it.
var pending = map[string]string{
	"serve": "#12",
	"gen":   "#12",
	"bench": "#22",
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "-h", "--help", "help":
		fmt.Fprint(stdout, usage)
		return 0
	case "--version", "version":
		fmt.Fprintln(stdout, "sensordash", version)
		return 0
	}
	if ticket, ok := pending[args[0]]; ok {
		fmt.Fprintf(stderr, "sensordash %s: not implemented yet: see %s\n", args[0], ticket)
		return 1
	}
	fmt.Fprintf(stderr, "sensordash: unknown command %q\n\n%s", args[0], usage)
	return 2
}
