#!/usr/bin/env python
# Cross compile Go source for various platforms. Based on https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04

import sys, os, subprocess

# sys/arch
platforms = [ "linux/amd64", "darwin/amd64", "linux/arm" ]

if len(sys.argv) != 2:
    print("Usage: {} <source.go>".format(sys.argv[0]))
    sys.exit(1)

source_name = sys.argv[1]
base = os.path.basename(source_name)
name = os.path.splitext(base)[0]

for platform in platforms:
    sys, arch = platform.split('/')
    bin_name = "{}-{}-{}".format(name, sys, arch)
    if sys == "windows":
        bin_name += ".exe"

    e = dict(os.environ)   # Make a copy of the current environment
    e['GOOS'] = sys
    e['GOARCH'] = arch
    subprocess.Popen(['go', 'build', '-o', bin_name, source_name], env=e)
