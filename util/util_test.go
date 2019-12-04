package util

import "testing"

func TestIsEmpty(t *testing.T) {
	type testpair struct {
		line    string
		isempty bool
	}

	tests := []testpair{
		{"", true},
		{" ", true},
		{"    ", true},
		{"\t", true},
		{"#", false},
		{"//", false},
		{"a", false},
		{" a", false},
	}

	for _, pair := range tests {
		v := IsEmpty(pair.line)
		if v != pair.isempty {
			t.Fatalf("For [%s] expected %v got %v\n", pair.line, pair.isempty, v)
		}
	}
}

func TestIsComment(t *testing.T) {
	type testpair struct {
		line      string
		isComment bool
	}

	tests := []testpair{
		// no comments
		{"", false},
		{"ls -l", false},
		{" ls -l", false},
		{"/usr/bin/perl -e 'print \"hello\n\"'", false},
		{"ls -l //etc/passwd", false},
		{"ls -l //etc/passwd # comment", false},
		{" ls -l /etc/passwd // comment", false},

		// bash-style comments
		{"#", true},
		{"##", true},
		{"###", true},
		{"# comment", true},
		{" # comment", true},
		{"#ls -l", true},
		{" #ls -l", true},
		{"#/usr/bin/perl -e 'print \"hello\n\"'", true},

		// golang-style comments
		{"// comment", true},
		{" // comment", true},
		{"//ls -l", true},
		{" //ls -l", true},
		{"///usr/bin/perl -e 'print \"hello\n\"'", true},
	}

	for _, pair := range tests {
		v := IsComment(pair.line)
		if v != pair.isComment {
			t.Fatalf("For [%s] expected %v got %v\n", pair.line, pair.isComment, v)
		}
	}
}
