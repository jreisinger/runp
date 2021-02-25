## About

`runp` is a simple command line tool that runs (shell) commands in parallel to save time. It's easy to install since it's a single binary. It's been tested on Linux (amd64 and arm) and MacOS/darwin (amd64).

```
# File containing commands we want to run in parallel.
$ cat cleanup.txt
kubectl delete all -l what=ckad
kubectl delete ing -l what=ckad
kubectl delete cm -l what=ckad
kubectl delete secret -l what=ckad

# Run commands from the file in parallel.
$ runp cleanup.txt 
--> OK (1.05s): /bin/sh -c "kubectl delete secret -l what=ckad"
No resources found
--> OK (1.05s): /bin/sh -c "kubectl delete ing -l what=ckad"
No resources found
--> OK (1.08s): /bin/sh -c "kubectl delete cm -l what=ckad"
No resources found
--> OK (3.78s): /bin/sh -c "kubectl delete all -l what=ckad"
No resources found
```

You might also like to see a related blog [post](https://reisinge.net/blog/2019-12-17-runp).

## Installation

Download the latest [release](https://github.com/jreisinger/runp/releases) to your `bin` folder (or some other folder on your `PATH`) and make it executable:

```
# Choose your system and architecture.
export SYS=linux  # or darwin
export ARCH=amd64 # or arm

# Download the binary and make it executable.
curl -L https://github.com/jreisinger/runp/releases/latest/download/runp-$SYS-$ARCH \
-o ~/bin/runp
chmod u+x ~/bin/runp
```

## Description and usage examples

Commands (or parts of them) can be read from files or stdin and must be separated by newlines. Comments and empty lines are ignored.

You can use shell variables in the commands.

`runp` exit status is 0 if all commands exit with 0 (OK).

`runp` prints a progress bar and info about command's execution (OK/ERR, run time, command to run) to stderr. Otherwise stdin and stderr works as you would expect. 

### Compress files

```
find . -iname '*.txt' | runp -p 'gzip --best'
```

We used `-p` to add prefix string to the commands (filenames in this case)

### Measure HTTP request + response time

```
export CURL="curl -w 'time_total:  %{time_total}\n' -o /dev/null -s "
for n in {1..10}; do echo $CURL https://yahoo.net; done | runp -q # 10 requests
```

### Find open TCP ports

```
cat << EOF > /tmp/host-port.txt
localhost 22
localhost 80
localhost 81
127.0.0.1 443
127.0.0.1 444
scanme.nmap.org 22
scanme.nmap.org 443
EOF

cat /tmp/host-port.txt | runp -q -p 'nc -v -w2 -z' 2>&1 | egrep '(succeeded!|open)$'
```

We used `-q` to suppress output from `runp` itself. Then we redirect stderr to stdout since netcat prints its messages to stderr. This way we can `grep` netcat's messages.

## Development

Test and install (to `~/go/bin/`):

```
make install
```

Test and build (cross-compile for multiple platforms):

```
# NOTE: don't forget to bump up the version in main.go!
make release
```

Working with GitHub releases:

```
# List existing tags
git tag

# Add new tag
git tag -a v2.1.3 -m "split code into packages, use modules, cleanup"
# Push tags to remote
git push origin --tags

# Add files from `make release` and publish the release on GitHub

# Delete tag
git tag -d v2.0.2
# Delete remote tag
git push --delete origin v2.0.2
```
