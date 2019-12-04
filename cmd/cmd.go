package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/jreisinger/runp/util"
)

// Command represents a command.
type Command struct {
	CmdString  string
	CmdToShow  string
	CmdToRun   *exec.Cmd
	StdoutCh   chan<- string
	StderrCh   chan<- string
	ExitCodeCh chan<- int8
	NoShell    bool
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

// Run runs a command and writes its stdout, stderr and exit code to corresponding channels.
func (c Command) Run() {
	stderr, err := c.CmdToRun.StderrPipe()
	if err != nil {
		c.StderrCh <- fmt.Sprintf("creating stderr pipe for %s: %s\n", c.CmdToShow, err)
		c.StdoutCh <- fmt.Sprintf("%s", "")
	}

	stdout, err := c.CmdToRun.StdoutPipe()
	if err != nil {
		c.StderrCh <- fmt.Sprintf("creating stdout pipe for %s: %s\n", c.CmdToShow, err)
		c.StdoutCh <- fmt.Sprintf("%s", "")
	}

	start := time.Now()

	if err := c.CmdToRun.Start(); err != nil {
		c.StderrCh <- fmt.Sprintf("starting command %s: %s\n", c.CmdToShow, err)
		c.StdoutCh <- fmt.Sprintf("%s", "")
		return
	}

	slurpErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		c.StderrCh <- fmt.Sprintf("slurping stderr of %s: %s\n", c.CmdToShow, err)
		c.StdoutCh <- fmt.Sprintf("%s", "")
	}

	slurpOut, err := ioutil.ReadAll(stdout)
	if err != nil {
		c.StderrCh <- fmt.Sprintf("slurping stdout of %s: %s\n", c.CmdToShow, err)
		c.StdoutCh <- fmt.Sprintf("%s", "")
	}

	secs := time.Since(start).Seconds()

	if err := c.CmdToRun.Wait(); err != nil {

		c.StderrCh <- fmt.Sprintf("\r--> ERR (%.2fs): %s\n%s\n%s", secs, c.CmdToShow, err, slurpErr)
		c.StdoutCh <- fmt.Sprintf("%s", "")

		// Did the command return a non-zero exit code?
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus := exitError.Sys().(syscall.WaitStatus)
			c.ExitCodeCh <- int8(waitStatus.ExitStatus())
		}

		return
	}

	secs = time.Since(start).Seconds()

	c.StderrCh <- fmt.Sprintf("\r--> OK (%.2fs): %s\n%s", secs, c.CmdToShow, slurpErr)
	c.StdoutCh <- fmt.Sprintf("%s", slurpOut)
	c.ExitCodeCh <- int8(0)
}

// ReadCommands returns commands to execute.
func ReadCommands(args []string) []string {

	if len(args) == 0 {
		return readCommandsFromStdin()
	}

	var cmds []string
	for _, arg := range args {
		c := readCommandsFromFile(arg)
		cmds = append(cmds, c...)
	}
	return cmds
}

func readCommandsFromStdin() []string {

	cmds, err := readCommands(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	return cmds
}

func readCommandsFromFile(fileName string) []string {

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return nil
	}
	defer file.Close()

	cmds, err := readCommands(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	return cmds
}

func readCommands(file *os.File) ([]string, error) {
	var cmds []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// skip comments
		if util.IsComment(line) {
			continue
		}

		// skip empty lines
		if util.IsEmpty(line) {
			continue
		}

		cmds = append(cmds, line)
	}
	return cmds, scanner.Err()
}
