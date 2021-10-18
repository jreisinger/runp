package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jreisinger/runp/cmd"
	"github.com/jreisinger/runp/util"
)

//
// Command line flags and usage message.
//

var noshell = flag.Bool("n", false, "don't invoke shell and don't expand env. vars")
var prefix = flag.String("p", "", "prefix to put in front of the commands")
var suffix = flag.String("s", "", "suffix to put behind the commands")
var goroutines = flag.Int("g", 10, "max number of commands (goroutines) to run in parallel")
var quiet = flag.Bool("q", false, "be quiet")
var version = flag.Bool("v", false, "print version")
var help = flag.Bool("h", false, "print help")

func init() {
	flag.Usage = func() {
		desc := `Run commands from file(s) or stdin in parallel. Commands must be separated by
newlines. Comments and empty lines are ignored. https://github.com/jreisinger/runp`
		fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] [file ...]\n", desc, os.Args[0])
		flag.PrintDefaults()
	}
}

//
// Main.
//

func main() { // main itself runs in a goroutine
	flag.Parse()

	if *version {
		fmt.Printf("runp %s\n", "v4.1.0")
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	cmds := cmd.ReadCommands(flag.Args())

	stderrChan := make(chan string)
	stdoutChan := make(chan string)
	exitCodeChan := make(chan int8)
	commandChan := make(chan *cmd.Command)

	// Simple workload balancer does not run more than *goroutines gouroutines in parallel.
	for i := 0; i < *goroutines; i++ {
		go cmd.Runner(commandChan)
	}

	if !*quiet {
		go util.ProgressBar()
	}

	for _, command := range cmds {
		if *prefix != "" {
			command = *prefix + " " + command
		}
		if *suffix != "" {
			command = command + " " + *suffix
		}
		c := cmd.Command{CmdString: command, StdoutCh: stdoutChan, StderrCh: stderrChan, ExitCodeCh: exitCodeChan, NoShell: *noshell, Quiet: *quiet}
		go cmd.Dispatcher(&c, commandChan)
	}

	var exitCodesSum int

	for range cmds {
		fmt.Fprint(os.Stderr, <-stderrChan)
		fmt.Fprint(os.Stdout, <-stdoutChan)
		exitCodesSum += int(<-exitCodeChan)
	}

	if exitCodesSum > 0 {
		os.Exit(1)
	}
}
