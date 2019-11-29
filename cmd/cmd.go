package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/jreisinger/runp/pkg/util"
)

// Command represents a command.
type Command struct {
	CmdString string
	CmdToShow string
	CmdToRun  *exec.Cmd
	StdoutCh  chan<- string
	StderrCh  chan<- string
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
		return
	}

	secs = time.Since(start).Seconds()

	c.StderrCh <- fmt.Sprintf("\r--> OK (%.2fs): %s\n%s", secs, c.CmdToShow, slurpErr)
	c.StdoutCh <- fmt.Sprintf("%s", slurpOut)
}

// ReadCommands reads in commands.
func ReadCommands(file *os.File) ([]string, error) {
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
