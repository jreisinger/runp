`runp` runs commands in parallel. It's useful when you want to run multiple commands (like those in `commands` folder) at the same time.

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
ln -sf runp-<sys>-<arch> runp
```

Run:

```
./runp <file-with-commands>
```

Example use case - installing my `vim` plugins (+ some other stuff):

```
$ ./runp commands/test.txt
--> OK (0.01s): /bin/sh -c "ls"
--> ERR (0.02s): /bin/sh -c "blah"
/bin/sh: 1: blah: not found
exit status 127
--> OK (3.02s): /bin/sh -c "sleep 3"
--> OK (5.02s): /bin/sh -c "sleep 5"
--> OK (9.01s): /bin/sh -c "sleep 9"
```

It took 9.01 seconds as opposed to the sum of all times as it would in case the commands run sequentially.
