`runp` -- run commands in parrallel. Useful when you want to run multiple
commands (like those in `commands` folder) at the same time.

Build for multiple systems and architectures:

```
./go-cross-compile.sh runp.go
ln -sf runp-<sys>-<arch> runp
```

Run:

```
./runp <file-with-commands>
```
