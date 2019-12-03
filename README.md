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

`runp` exit status is 0 if all commands exit with 0 (OK). stdin and stderr works as usual.

### Ping several hosts (read from stdin)

```
$ runp -p 'ping -c 2 -W 2' > /dev/null
localhost
1.1.1.1 # Clouflare
8.8.8.8 # Google
--> OK (1.06s): /bin/sh -c "ping -c 2 -W 2 localhost"
--> OK (1.07s): /bin/sh -c "ping -c 2 -W 2 1.1.1.1 # Clouflare"
--> OK (1.07s): /bin/sh -c "ping -c 2 -W 2 8.8.8.8 # Google"
```

Press `Ctrl-D` when done entering the host names. Running all the commands took only 1.07 second as opposed to the sum of all times. We suppressed the printing of commands' stdout by redirecting stdout to `/dev/null`.

### Get directories' sizes (read from stdin)

```
$ echo -e "/home\n/etc\n/tmp\n/data/backup\n/data/public" | sudo runp -n -p 'du -sh' 2> /dev/null 
4.7M	/tmp
7.1M	/etc
943M	/home
416G	/data/public
292G	/data/backup
```

We suppressed the printing of progress bar and info about command's execution (OK/ERR, run time, command to run) by discarding stderr.

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
