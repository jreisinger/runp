## About

`runp` runs shell commands in parallel. It's useful when you want to run multiple commands (like those in `commands` folder) at the same time.

Installation:

```
export SYS=linux  # darwin
export ARCH=amd64 # arm
curl --location https://github.com/jreisinger/runp/releases/latest/download/runp-$SYS-$ARCH \
--output ~/bin/runp && chmod u+x ~/bin/runp
```

Usage:

```
$ runp [-v] <file-with-commands>
```

Example:

```
$ runp commands/test.txt
--> OK (0.01s): /bin/sh -c "ls"
--> ERR (0.02s): /bin/sh -c "blah"
/bin/sh: 1: blah: not found
exit status 127
--> OK (3.02s): /bin/sh -c "sleep 3"
--> OK (5.02s): /bin/sh -c "sleep 5"
--> OK (9.01s): /bin/sh -c "sleep 9"
```

It took 9.01 seconds as opposed to the sum of all times as it would in case the commands run sequentially. If the command exits with 0 runp prints `OK`. Otherwise it prints `ERR` (in red) and STDERR. If you want to see also STDOUT use the `-v` switch.

## Development

Prep:

```
export GOPATH=`pwd`
go get -u github.com/fatih/color
```

Test:

```
go test
```

Build (for multiple systems and architectures):

```
./go-cross-compile.py runp.go
ln -sf runp-$SYS-$ARCH runp
./runp
```
