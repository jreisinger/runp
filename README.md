## About

`runp` is a command line tool that runs (shell) commands in parallel or concurrently. It's useful when you want to run multiple commands at once to save time. It's somehow similar to the GNU [parallel](https://www.gnu.org/software/parallel/) tool.

## Installation

Download the latest [release](https://github.com/jreisinger/runp/releases) to your `bin` folder (or some other folder on your `PATH`) and make it executable:

```
export SYS=linux  # darwin
export ARCH=amd64 # arm
curl -L https://github.com/jreisinger/runp/releases/latest/download/runp-$SYS-$ARCH -o ~/bin/runp
chmod u+x ~/bin/runp
```

## Usage examples

You can use shell variables in the commands. Commands have to be separated by newlines. Empty lines and comments are ignored.

### Run some test commands (read from a file)

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

Running all the commands took only 9.02 seconds as opposed to the sum of all times.

### Get directories' sizes (read from stdin)

```
$ echo -e "/home\n/etc\n/tmp\n/data/backup\n/data/public" | sudo runp -n -p 'du -sh' 2> /dev/null 
4.7M	/tmp
7.1M	/etc
943M	/home
416G	/data/public
292G	/data/backup
```

We surpressed the printing of progress bar and info about command's execution (OK/ERR, run time, command to run) by discarding stderr.

### Get Jupiter images from NASA

```
base='https://images-api.nasa.gov/search'
query='jupiter'
desc='planet'
type='image'
curl -s "$base?q=$query&description=$desc&media_type=$type" | \
jq -r .collection.items[].href | head -50 | runp -p 'curl -s' | jq -r .[] | grep large | \
runp -p 'curl -s -L -O'
```

### Compress images

```
find -iname '*.jpg' | runp -p 'gzip --best'
```

## Development

Test:

```
make test
```

Test and install:

```
make install
```

Test and build (cross-compile for multiple platforms):

```
make release
```
