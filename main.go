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
		fmt.Printf("runp %s\n", "v2.0.2")
		os.Exit(0)
	}

	var cmds []string // commands to execute

	if len(flag.Args()) == 0 { // get commands to execute from STDIN
		fileCmds, err := cmd.ReadCommands(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		cmds = append(cmds, fileCmds...)
	} else {
		for _, arg := range flag.Args() { // get commands to execute from file(s)
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