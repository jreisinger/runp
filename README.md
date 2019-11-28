## About

`runp` runs shell commands in parallel (or concurrently). It's useful when you want to run multiple shell commands at once to save time. It's somehow similar to the GNU [parallel](https://www.gnu.org/software/parallel/) tool.

## Installation

Download the latest [release](https://github.com/jreisinger/runp/releases) to your `bin` folder and make it executable:

```
export SYS=linux  # darwin
export ARCH=amd64 # arm
curl --location https://github.com/jreisinger/runp/releases/latest/download/runp-$SYS-$ARCH \
--output ~/bin/runp && chmod u+x ~/bin/runp
```

## Usage examples

You can use shell variables in the commands. Commands have to be separated by newlines. Empty lines and comments are ignored.

### Run test commands (file)

```
$ runp commands/test.txt > /dev/null
--> OK (0.01s): /bin/sh -c "ls /Users/reisinge/github/runp # 'PWD' shell variable is used here"
--> ERR (0.01s): /bin/sh -c "blah"
exit status 127
/bin/sh: blah: command not found
--> OK (3.02s): /bin/sh -c "sleep 3"
--> OK (5.02s): /bin/sh -c "sleep 5"
--> OK (9.02s): /bin/sh -c "sleep 9"
```

Running all the commands took 9.02 seconds. As opposed to the sum of all times in case the commands ran sequentially.

### Get directories' sizes (stdin)

```
$ echo -e "/home\n/etc\n/tmp\n/data/backup\n/data/public" | sudo runp -p 'du -sh'
--> OK (0.04s): /bin/sh -c "du -sh /tmp"
472K    /tmp
--> OK (0.09s): /bin/sh -c "du -sh /etc"
7.1M    /etc
--> OK (0.49s): /bin/sh -c "du -sh /home"
933M    /home
--> OK (5.59s): /bin/sh -c "du -sh /data/backup"
292G    /data/backup
--> OK (32.45s): /bin/sh -c "du -sh /data/public"
415G    /data/public
```

### Get some NASA images (stdin)

```
base='https://images-api.nasa.gov/search'
query='jupiter'
desc='planet'
type='image'
curl -s "$base?q=$query&description=$desc&media_type=$type" | \
jq -r .collection.items[].href | head -50 | runp -p 'curl -s' | jq -r .[] | grep large | \
runp -p 'curl -s -L -O'
```

## Development

Prep:

```
export GOPATH=`pwd`
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
