`runp` - run commands in parallel. Useful when you want to run multiple
commands (like those in `commands` folder) at the same time.

Build (for multiple systems and architectures):

```
./go-cross-compile.py runp.go
ln -sf runp-<sys>-<arch> runp
```

Run:

```
./runp <file-with-commands>
```

Example use case - installing my `vim` plugins:

```
$ ./runp commands/install-my-stuff.txt
--> OK (0.40s): /bin/sh -c "curl -L -o /Users/reisinge/.git-completion.bash https://raw.githubusercontent.com/git/git/master/contrib/completion/git-completion.bash"
--> OK (0.46s): /bin/sh -c "curl -L -o /Users/reisinge/.vim/autoload/pathogen.vim https://raw.github.com/tpope/vim-pathogen/master/autoload/pathogen.vim"
--> OK (1.13s): /bin/sh -c "rm -rf /Users/reisinge/.vim/bundle/ansible-vim && git clone https://github.com/pearofducks/ansible-vim.git /Users/reisinge/.vim/bundle/ansible-vim"
--> OK (1.17s): /bin/sh -c "rm -rf /Users/reisinge/.vim/bundle/grep.vim && git clone https://github.com/yegappan/grep.git /Users/reisinge/.vim/bundle/grep.vim"
--> OK (1.27s): /bin/sh -c "rm -rf /Users/reisinge/.vim/bundle/vim-nerdtree-tabs && git clone https://github.com/jistr/vim-nerdtree-tabs.git /Users/reisinge/.vim/bundle/vim-nerdtree-tabs"
--> OK (1.31s): /bin/sh -c "rm -rf /Users/reisinge/.vim/bundle/bufexplorer && git clone https://github.com/jlanzarotta/bufexplorer.git /Users/reisinge/.vim/bundle/bufexplorer"
--> OK (1.53s): /bin/sh -c "rm -rf /Users/reisinge/.vim/bundle/vim-markdown && git clone https://github.com/plasticboy/vim-markdown.git /Users/reisinge/.vim/bundle/vim-markdown"
--> OK (2.26s): /bin/sh -c "rm -rf /Users/reisinge/.vim/bundle/nerdtree && git clone https://github.com/scrooloose/nerdtree.git /Users/reisinge/.vim/bundle/nerdtree"
--> OK (2.33s): /bin/sh -c "rm -rf /Users/reisinge/.vim/pack/dist/start/vim-airline && git clone https://github.com/vim-airline/vim-airline /Users/reisinge/.vim/pack/dist/start/vim-airline"
```

It took 2.33 seconds as opposed to the sum of all times as it would in case the
commands run sequentially.
