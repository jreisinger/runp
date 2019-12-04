package cmd

import (
	"testing"
)

func TestCommandPrepare(t *testing.T) {
	noshell := false
	stderrChan := make(chan string)
	stdoutChan := make(chan string)

	// empty command
	cmd := ""
	c := Command{CmdString: cmd, StdoutCh: stdoutChan, StderrCh: stderrChan, NoShell: noshell}
	c.Prepare()
	if c.CmdString != "" {
		t.Fatalf("CmdString is not empty: %v", c.CmdString)
	}
	if c.CmdToShow != "/bin/sh -c \"\"" {
		t.Fatalf("CmdToShow is not empty: %v", c.CmdToShow)
	}

	// basic command
	cmd = "ls -l"
	c = Command{CmdString: cmd, StdoutCh: stdoutChan, StderrCh: stderrChan, NoShell: noshell}
	c.Prepare()
	if c.CmdString != "ls -l" {
		t.Fatalf("CmdString is wrong")
	}
	if c.CmdToShow != "/bin/sh -c \"ls -l\"" {
		t.Fatalf("CmdToShow is wrong")
	}
}
