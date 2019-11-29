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
$ echo -e "/home\n/etc\n/tmp\n/data/backup\n/data/public" | sudo runp -n -p 'du -sh'
--> OK (0.02s): du -sh /tmp
476K	/tmp
--> OK (0.03s): du -sh /etc
7.1M	/etc
--> OK (0.33s): du -sh /home
933M	/home
--> OK (0.74s): du -sh /data/public
416G	/data/public
--> OK (2.55s): du -sh /data/backup
292G	/data/backup
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

### Compress images

```
find -iname '*.jpg' | runp -p 'gzip --best'
```

## Development

Test and build (cross-compile):

```
make release
```

Install:

```
make install
```
