## About

`runp` is a simple command line tool that runs (shell) commands in parallel. It's useful when you want to run multiple commands at once to save time. It's easy to install since it's a single binary. It's been tested on Linux (amd64 and arm) and MacOS/darwin (amd64).

There's also a related blog [post](https://jreisinger.github.io/blog2/posts/runp/) (with a movie! :-).

## Installation

Download the latest [release](https://github.com/jreisinger/runp/releases) to your `bin` folder (or some other folder on your `PATH`) and make it executable:

```
# Choose your system and architecture.
export SYS=linux  # or darwin
export ARCH=amd64 # or arm

# Download the binary and make it executable.
curl -L https://github.com/jreisinger/runp/releases/latest/download/runp-$SYS-$ARCH -o ~/bin/runp
chmod u+x ~/bin/runp
```

## Description and Usage examples

Commands (or parts of them) can be read from files or stdin and must be separated by newlines. Comments and empty lines are ignored.

You can use shell variables in the commands.

`runp` exit status is 0 if all commands exit with 0 (OK).

`runp` prints a progress bar and info about command's execution (OK/ERR, run time, command to run) to stderr. Otherwise stdin and stderr works as you would expect. 

### Run some test commands (read from file)

```
# Create a file containing commands to run in parallel.
cat << EOF > /tmp/test-commands.txt
sleep 5
sleep 3
blah     # this will fail
ls $PWD  # PWD shell variable is used here
EOF

# Run commands from the file.
runp /tmp/test-commands.txt > /dev/null
```

We suppressed the printing of commands' stdout by redirecting it to `/dev/null`.

### Ping several hosts and see packet loss (read from stdin)

```
runp -p 'ping -c 5 -W 2' -s '| grep loss' # first copy this line then the rest
localhost
1.1.1.1
8.8.8.8
# Press `Ctrl-D` when done entering the hosts
```

We used `-p` and `-s` to add prefix and suffix strings to the commands (hosts in this case).

### Get directories' sizes (read from stdin)

```
echo -e "$HOME\n/etc\n/tmp" | runp -q -p 'du -sh'
```

We suppressed the printing of `runp`'s progress bar and info about command's execution by using the `-q` flag.

### Get Jupiter images from NASA

```
base='https://images-api.nasa.gov/search'
query='jupiter'
desc='planet'
type='image'
curl -s "$base?q=$query&description=$desc&media_type=$type" | \
jq -r .collection.items[].href | runp -p 'curl -s' | jq -r .[] | grep large | \
runp -p 'curl -s -L -O'
```

If you want to see how much time you just saved add `-g 1` to both `runp` invocations :-). It will execute the commands sequentially.

### Compress images

```
find . -iname '*.jpg' | runp -p 'gzip --best'
```

### Measure HTTP request + response time

```
export CURL="curl -w 'time_total:  %{time_total}\n' -o /dev/null -s https://golang.org/"
perl -wE 'for (1..10) { say $ENV{CURL} }' | runp -q # make 10 requests
```

### Find open TCP ports

```
$ cat host-port.txt
localhost 80
localhost 81
127.0.0.1 443
127.0.0.1 444
localhost 22
$ cat host-port.txt | runp -p 'netcat -v -w2 -z' -q 2>&1 | egrep 'open$'
localhost [127.0.0.1] 443 (https) open
localhost [127.0.0.1] 80 (http) open
localhost [127.0.0.1] 22 (ssh) open
```

We used `-q` to suppress output from `runp` itself. Then we redirect stderr to stdout since netcat prints its messages to stderr. This way web can `grep` netcat's messages.

## Development

Test and install (to `~/go/bin/`):

```
make install
```

Test and build (cross-compile for multiple platforms):

```
make release
```

Working with GitHub releases:

*NOTE*: don't forget to bump up the version in `main.go`!

```
# List existing tags
git tag

# Add new tag
git tag -a v2.1.3 -m "split code into packages, use modules, cleanup"
# Push tags to remote
git push origin --tags

# Delete tag
git tag -d v2.0.2
# Delete remote tag
git push --delete origin v2.0.2
```
