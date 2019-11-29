package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jreisinger/runp/cmd"
	"github.com/jreisinger/runp/pkg/util"
)

func main() { // main itself runs in a goroutine
	// Usage and command line args.

	flag.Usage = usage

	noshell := flag.Bool("n", false, "don't invoke shell and don't expand env. vars")
	version := flag.Bool("V", false, "print version")
	prefix := flag.String("p", "", "prefix to put in front of the commands")
	suffix := flag.String("s", "", "suffix to put behind the commands")

	flag.Parse()

	if *version {
		fmt.Printf("runp %s\n", "v2.0.1")
		os.Exit(0)
	}

	// all commands to execute
	var cmds []string

	if len(flag.Args()) == 0 {
		// Get commands to execute from STDIN.
		fileCmds, err := cmd.ReadCommands(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		cmds = append(cmds, fileCmds...)
	} else {
		// Get commands to execute from files.
		for _, arg := range flag.Args() {
			file, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				continue
			}
			defer file.Close()

			fileCmds, err := cmd.ReadCommands(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			cmds = append(cmds, fileCmds...)
		}
	}

	// Run commands in parallel.

	stderrChan := make(chan string)
	stdoutChan := make(chan string)

	go util.ProgressBar()
	for _, command := range cmds {
		if *prefix != "" {
			command = *prefix + " " + command
		}
		if *suffix != "" {
			command = command + " " + *suffix
		}
		c := cmd.Command{CmdString: command, StdoutCh: stdoutChan, StderrCh: stderrChan, NoShell: *noshell}
		c.Prepare()
		go c.Run()
	}

	for range cmds {
		fmt.Fprint(os.Stderr, <-stderrChan)
		fmt.Fprint(os.Stdout, <-stdoutChan)
	}

}

func usage() {
	desc := `Run commands from file(s) or stdin in parallel. Commands are
separated by newlines. Comments and empty lines are skipped.`
	fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] [file ...]\n", desc, os.Args[0])
	flag.PrintDefaults()
}
