package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

func main() { // main itself runs in a goroutine
	// Usage and command line args.

	flag.Usage = usage

	verbose := flag.Bool("v", false, "be verbose")
	noshell := flag.Bool("n", false, "don't invoke shell and don't expand env. vars")
	version := flag.Bool("V", false, "print version")
	prefix := flag.String("p", "", "prefix to put in front of the commands")
	suffix := flag.String("s", "", "suffix to put behind the commands")

	flag.Parse()

	if *version {
		fmt.Printf("runp %s\n", "v1.2.0")
		os.Exit(0)
	}

	// All commands to execute.
	var cmds []string

	if len(flag.Args()) == 0 {
		// Get commands to execute from STDIN.
		fileCmds, err := readCommands(os.Stdin)
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

			fileCmds, err := readCommands(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v\n", err)
				os.Exit(1)
			}
			cmds = append(cmds, fileCmds...)
		}
	}

	// Run commands in parallel.

	ch := make(chan string)

	go progressBar()
	for _, cmd := range cmds {
		if *prefix != "" {
			cmd = *prefix + " " + cmd
		}
		if *suffix != "" {
			cmd = cmd + " " + *suffix
		}
		c := Command{CmdString: cmd, Channel: ch, Verbose: *verbose, NoShell: *noshell}
		c.Prepare()
		go c.Run()
	}

	for range cmds {
		fmt.Print(<-ch) // receive from channel ch
	}
}

func progressBar() {
	for {
		count := 0
		for {
			count++
			time.Sleep(100 * time.Millisecond)
			if count == 3 {
				fmt.Printf(">\r")
				count = 0
				continue
			}
			fmt.Printf("-")
		}
	}
}

func usage() {
	desc := `Run commands from file(s) or stdin in parallel. By default, shell is
invoked and env. vars are expanded. Comments and empty lines are skipped.`
	fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] [file ...]\n", desc, os.Args[0])
	flag.PrintDefaults()
}

// Command represents a command.
type Command struct {
	CmdString string
	CmdToShow string
	CmdToRun  *exec.Cmd
	Channel   chan<- string
	Verbose   bool
	NoShell   bool
}

// Prepare prepares a command to be run.
func (c *Command) Prepare() {
	if c.NoShell {
		parts := strings.Split(c.CmdString, " ")
		c.CmdToRun = exec.Command(parts[0], parts[1:]...)
		c.CmdToShow = c.CmdString
	} else {
		c.CmdString = os.ExpandEnv(c.CmdString) // expand ${var} or $var
		shellToUse := "/bin/sh"
		c.CmdToRun = exec.Command(shellToUse, "-c", c.CmdString)
		c.CmdToShow = shellToUse + " -c " + strconv.Quote(strings.Join(c.CmdToRun.Args[2:], " "))
	}
}

// Run runs a command.
func (c Command) Run() {
	start := time.Now()
	stdoutStderr, err := c.CmdToRun.CombinedOutput()
	secs := time.Since(start).Seconds()
	if err != nil {
		red := color.New(color.FgRed).SprintFunc()
		//c.Channel <- fmt.Sprintf("\r--> ERR (%.2fs): %s\n%s%s\n", secs, c.CmdToShow, stdoutStderr, err)
		c.Channel <- fmt.Sprintf("\r--> %s (%.2fs): %s\n%s%s\n", red("ERR"), secs, c.CmdToShow, stdoutStderr, err)
		return
	}

	if c.Verbose {
		c.Channel <- fmt.Sprintf("\r--> OK (%.2fs): %s\n%s", secs, c.CmdToShow, stdoutStderr)
	} else {
		c.Channel <- fmt.Sprintf("\r--> OK (%.2fs): %s\n", secs, c.CmdToShow)
	}
}

// readCommands reads in commands.
func readCommands(file *os.File) ([]string, error) {
	var cmds []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// skip comments
		if isComment(line) {
			continue
		}

		// skip empty lines
		if isEmpty(line) {
			continue
		}

		cmds = append(cmds, line)
	}
	return cmds, scanner.Err()
}

// isEmpty returns true if line is empty.
func isEmpty(line string) bool {
	var emptyLine = regexp.MustCompile(`^\s*$`)
	return emptyLine.MatchString(line)
}

// isComment returns true if line is a comment.
func isComment(line string) bool {
	match, _ := regexp.MatchString(`^\s*(#|//)`, line)
	return match
}
