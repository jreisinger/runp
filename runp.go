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

	flag.Parse()

	if *version {
		fmt.Printf("runp %s\n", "v1.0.0")
		os.Exit(0)
	}

	if len(flag.Args()) != 1 {
		usage()
		os.Exit(1)
	}

	// Get commands to execute from a file.

	var cmds []string

	cmds, err := readCommands(flag.Args()[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// Run commands in parallel.

	ch := make(chan string)

	go progressBar()
	for _, cmd := range cmds {
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
	desc := `Run commands defined in a file in parallel. By default, shell is
invoked and env. vars are expanded. Comments and empty lines are skipped.`
	fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] commands.txt\n", desc, os.Args[0])
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

// readCommands reads command strings from a file.
func readCommands(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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
