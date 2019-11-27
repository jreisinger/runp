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

### Test commands

```
$ runp commands/test.txt
--> ERR (0.00s): /bin/sh -c "blah"
/bin/sh: 1: blah: not found
exit status 127
--> OK (0.00s): /bin/sh -c "ls /home/reisinge/github/runp # 'PWD' shell variable is used here"
--> OK (3.00s): /bin/sh -c "sleep 3"
--> OK (5.00s): /bin/sh -c "sleep 5"
--> OK (9.01s): /bin/sh -c "sleep 9"
```

Running all the commands took 9.01 seconds. As opposed to the sum of all times in case the commands ran sequentially. If the command exits with 0 runp prints `OK`. Otherwise it prints `ERR` (in red) and STDERR. If you want to see also STDOUT use the `-v` switch.

### Get some NASA images in parallel

```
base='https://images-api.nasa.gov/search'
query='apollo%2011'
desc='moon%20landing'
type='image'
for url in $(curl "$base?q=$query&description=$desc&media_type=$type" | \
jq -r .collection.items[].href | head); do curl $url | jq -r .[] | grep large; done 2> /dev/null | \
runp -p 'curl -L -O'
```

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
