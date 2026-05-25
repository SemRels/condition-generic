package main

import (
	"fmt"
	"io"
	"os"

	plugin "github.com/SemRels/condition-generic/internal/plugin"
)

func run(getenv func(string) string, stderr io.Writer) int {
	c := plugin.NewWithEnv(getenv)
	if err := c.Check(); err != nil {
		fmt.Fprintln(stderr, "condition-generic:", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Getenv, os.Stderr))
}
