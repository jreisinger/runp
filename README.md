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

Commands can be read from file(s) or stdin and must be separated by newlines. Comments and empty lines are ignored.

You can use shell variables in the commands. `runp` exit status is 0 if all commands exit with 0 (OK). stdin and stderr work as usual. 

### Run some test commands (read from file)

```
cat << EOF > /tmp/commands.txt
sleep 5
sleep 3
blah     # this will fail
ls $PWD  # 'PWD' shell variable is used here
EOF

runp /tmp/commands.txt > /dev/null
```

### Get directories' sizes (read from stdin)

```
echo -e "$HOME\n/etc\n/tmp" | runp -n -p 'du -sh' 2> /dev/null 
```

We suppressed the printing of progress bar and info about command's execution (OK/ERR, run time, command to run) by redirecting stderr to `/dev/null`.

### Ping several hosts and see packet loss (read from stdin)

```
runp -p 'ping -c 5 -W 2' -s '| grep loss'
localhost
1.1.1.1
8.8.8.8
# Press `Ctrl-D` when done entering the hosts
```

We used `-p` and `-s` to add prefix and suffix strings to the commands (hosts in this case).

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
find . -iname '*.jpg' | runp -p 'gzip --best'
```

### Measure HTTP request and response time

```
export CURL="curl -w 'time_total:  %{time_total}\n' -o /dev/null -s https://golang.org/"
perl -E 'for (1..100) { say $ENV{CURL} }' | runp 2> /dev/null
```

## Development

Test and install:

```
make install
```

Test and build (cross-compile for multiple platforms):

```
make release
```
