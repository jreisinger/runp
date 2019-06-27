package main

import (
    "testing"
)

func TestCommandPrepare(t *testing.T) {
    verbose := false
    noshell := false
    ch := make(chan string)

    // empty command
    cmd := ""
    c := Command{CmdString: cmd, Channel: ch, Verbose: verbose, NoShell: noshell}
    c.Prepare()
    if c.CmdString != "" {
        t.Fatalf("CmdString is not empty: %v", c.CmdString)
    }
    if c.CmdToShow != "/bin/sh -c \"\"" {
        t.Fatalf("CmdToShow is not empty: %v", c.CmdToShow)
    }

    // basic command
    cmd = "ls -l"
    c = Command{CmdString: cmd, Channel: ch, Verbose: verbose, NoShell: noshell}
    c.Prepare()
    if c.CmdString != "ls -l" {
        t.Fatalf("CmdString is wrong")
    }
    if c.CmdToShow != "/bin/sh -c \"ls -l\"" {
        t.Fatalf("CmdToShow is wrong")
    }
}

func TestIsComment(t *testing.T) {
    type testpair struct {
        line        string
        isComment   bool
    }

    tests := []testpair{
        // no comments
        { "", false },
        { "ls -l", false },
        { "/urs/bin/perl -e 'print \"hello\n\"'", false },
        // comments
        { "// comment", true },
        { "# comment", true },
    }

    for _, pair := range tests {
        v := isComment(pair.line)
        if v != pair.isComment {
            t.Fatal(
                "For", pair.line,
                "expected", pair.isComment,
                "got", v,
            )
        }
    }
}