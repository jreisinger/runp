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

Commands can be read from stdin or from file(s) and must be separated by newlines. Comments and empty lines are ignored.

You can use shell variables in the commands. `runp` exit status is 0 if all commands exit with 0 (OK). stdin and stderr work as usual. 

### Ping several hosts (read from stdin)

```
runp -p 'ping -c 2 -W 2' > /dev/null
localhost
1.1.1.1 # Clouflare
8.8.8.8
# Press `Ctrl-D` when done entering the host names.

```

We suppressed the printing of commands' stdout by redirecting stdout to `/dev/null`.

### Get directories' sizes (read from stdin)

```
echo -e "$HOME\n/etc\n/tmp" | sudo runp -n -p 'du -sh' 2> /dev/null 
```

We suppressed the printing of progress bar and info about command's execution (OK/ERR, run time, command to run) by redirecting stderr to `/dev/null`.

### Run some test commands (read from file)

```
cat << EOF > /tmp/commands.txt
sleep 3
sleep 5

blah

# comment
sleep 9

ls $PWD # 'PWD' shell variable is used here
EOF

runp /tmp/commands.txt > /dev/null
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

Test and install:

```
make install
```

Test and build (cross-compile for multiple platforms):

```
make release
```
