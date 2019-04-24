package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func usage() {
	desc := `Run commands defined in a file in parallel. By default, shell is invoked and
env. vars are expanded. Source: https://github.com/jreisinger/runp`
	fmt.Fprintf(os.Stderr, "%s\n\nUsage: %s [options] commands.txt\n", desc, os.Args[0])
	flag.PrintDefaults()
}

func main() { // main runs in a goroutine
	flag.Usage = usage

	verbose := flag.Bool("v", false, "be verbose")
	noshell := flag.Bool("n", false, "don't invoke shell and don't expand env. vars")

	flag.Parse()

	if len(flag.Args()) != 1 {
		usage()
		os.Exit(1)
	}

	log.SetPrefix("runp: ")
	log.SetFlags(0) // no extra info in log messages

	// Get commands to execute from a file.
	cmds, err := readCommands(flag.Args()[0])
	if err != nil {
		log.Fatalf("Error reading commands: %s. Exiting ...\n", err)
	}

	ch := make(chan string)

	for _, cmd := range cmds {
		go run(cmd, ch, verbose, noshell)
	}

	for range cmds {
		// receive from channel ch
		fmt.Print(<-ch)
	}
}

func readCommands(filePath string) ([]string, error) {
	// Open the file containing commands.
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
		match, _ := regexp.MatchString("^(#|/)", line)
		if match {
			continue
		}

		cmds = append(cmds, line)
	}
	return cmds, scanner.Err()
}

func run(command string, ch chan<- string, verbose *bool, noshell *bool) {
	var cmd *exec.Cmd
	var cmdToShow string

	if *noshell {
		parts := strings.Split(command, " ")
		cmd = exec.Command(parts[0], parts[1:]...)
		cmdToShow = strings.Join(cmd.Args, " ")
	} else {
		command = os.ExpandEnv(command) // expand ${var} or $var
		shellToUse := "/bin/sh"
		cmd = exec.Command(shellToUse, "-c", command)
		cmdToShow = shellToUse + " -c " + strconv.Quote(strings.Join(cmd.Args[2:], " "))
	}

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		ch <- fmt.Sprintf("--> ERR: %s\n%s%s\n", cmdToShow, stdoutStderr, err)
		return
	}

	if *verbose {
		ch <- fmt.Sprintf("--> OK: %s\n%s\n", cmdToShow, stdoutStderr)
	} else {
		ch <- fmt.Sprintf("--> OK: %s\n", cmdToShow)
	}
}
