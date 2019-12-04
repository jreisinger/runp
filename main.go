package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jreisinger/runp/cmd"
	"github.com/jreisinger/runp/pkg/util"
)

func usage() {
	desc := `Run commands from file(s) or stdin in parallel. Commands must be separated by
newlines. Comments and empty lines are ignored.`
	fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] [file ...]\n", desc, os.Args[0])
	flag.PrintDefaults()
}

func main() { // main itself runs in a goroutine
	flag.Usage = usage

	noshell := flag.Bool("n", false, "don't invoke shell and don't expand env. vars")
	version := flag.Bool("V", false, "print version")
	prefix := flag.String("p", "", "prefix to put in front of the commands")
	suffix := flag.String("s", "", "suffix to put behind the commands")

	flag.Parse()

	if *version {
		fmt.Printf("runp %s\n", "v2.1.3")
		os.Exit(0)
	}

	cmds := cmd.ReadCommands(flag.Args())

	stderrChan := make(chan string)
	stdoutChan := make(chan string)
	exitCodeChan := make(chan int8)

	go util.ProgressBar()
	for _, command := range cmds {
		if *prefix != "" {
			command = *prefix + " " + command
		}
		if *suffix != "" {
			command = command + " " + *suffix
		}
		c := cmd.Command{CmdString: command, StdoutCh: stdoutChan, StderrCh: stderrChan, ExitCodeCh: exitCodeChan, NoShell: *noshell}
		c.Prepare()
		go c.Run()
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
