# Testing

There are two "parts" of completion:

 - the selection parts - tested in `app_completion_ext_test.go`
 - the zsh/bash/fish integration parts - tested manually.

 ## `zsh` Manual E2E Tests

ChatGPT offers ways to unit test zsh completion, and there's a few examples on StackOverflow, but I'll just manually test for now till I feel up to trying to automate it.

Setup:

```zsh
# install the app to $GOBIN
go install ./completion/internal/testappcmd

cd ~/go/bin  # default $GOBIN

# install the completion to something on $FPATH
./testappcmd completion zsh > ~/fbin/_testappcmd

# open a new shell to load completions
./testappcmd ...  # tab away!!

# test tab completions
./testappcmd manual <TAB>
```

If it's not loading for some reason, remove `~/.zcompdump` and open a new shell to regenerate it.

## `bash` Manual E2E Tests

Setup:

```bash
# install the app to $GOBIN
go install ./completion/internal/testappcmd

# not my default shell, so switch to it
bash

cd ~/go/bin  # default $GOBIN

# source the completion script in your shell
source <(./testappcmd completion bash)

# test tab completions
./testappcmd manual <TAB>
```

## `fish` Manual E2E Tests

Setup:

```bash
# install the app to $GOBIN
go install ./completion/internal/testappcmd

# not my default shell, so switch to it
fish

cd ~/go/bin  # default $GOBIN

# fish needs to add it to the PATH
fish_add_path $PWD

# source the completion script in your shell
./testappcmd completion fish | source

# test tab completions
./testappcmd manual <TAB>
```

