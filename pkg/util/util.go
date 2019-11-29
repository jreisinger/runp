package util

import (
	"fmt"
	"os"
	"regexp"
	"time"
)

// ProgressBar draws a moving --> to show something is hapenning.
func ProgressBar() {
	for {
		count := 0
		for {
			count++
			time.Sleep(100 * time.Millisecond)
			if count == 3 {
				fmt.Fprintf(os.Stderr, ">\r")
				count = 0
				continue
			}
			fmt.Fprintf(os.Stderr, "-")
		}
	}
}

// IsEmpty returns true if line is empty.
func IsEmpty(line string) bool {
	var emptyLine = regexp.MustCompile(`^\s*$`)
	return emptyLine.MatchString(line)
}

// IsComment returns true if line is a comment.
func IsComment(line string) bool {
	match, _ := regexp.MatchString(`^\s*(#|//)`, line)
	return match
}
