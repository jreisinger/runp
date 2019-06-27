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
        { " ls -l", false },
        { "/usr/bin/perl -e 'print \"hello\n\"'", false },
        { "ls -l //etc/passwd", false },
        { "ls -l //etc/passwd # comment", false },
        { " ls -l /etc/passwd // comment", false },

        // bash-style comments
        { "#", true },
        { "##", true },
        { "###", true },
        { "# comment", true },
        { " # comment", true },
        { "#ls -l", true },
        { " #ls -l", true },
        { "#/usr/bin/perl -e 'print \"hello\n\"'", true },

        // golang-style comments
        { "// comment", true },
        { " // comment", true },
        { "//ls -l", true },
        { " //ls -l", true },
        { "///usr/bin/perl -e 'print \"hello\n\"'", true },
    }

    for _, pair := range tests {
        v := isComment(pair.line)
        if v != pair.isComment {
            t.Fatalf("For [%s] expected %v got %v\n", pair.line, pair.isComment, v)
        }
    }
}