# Testing

There are two "parts" of completion:

 - the selection parts - tested in `app_completion_ext_test.go`
 - the zsh integration parts - tested manually.

 ## Test `zsh` integration.

ChatGPT offers ways to unit test zsh completion, and there's a few examples on StackOverflow, but I'll just manually test for now till I feel up to trying to automate it.

```zsh
# install the app to $GOBIN
go install ./completion/internal/testappcmd
cd ~/go/bin  # default $GOBIN
# install the completion to something on $FPATH
./testappcmd --completion-script-zsh > ~/fbin/_testappcmd
# open a new shell to load completions
./testappcmd ...  # tab away!!
```

# Manual tests

Just tab after all of these...

TODO: organize these by checking the branches of the completion function

## `COMPLETION_TYPE_DIRECTORIES`

## `COMPLETION_TYPE_DIRECTORIES_FILES`

## `COMPLETION_TYPE_NONE`

## `COMPLETION_TYPE_VALUES`

./

## `COMPLETION_TYPE_VALUES_DESCRIPTIONS`

```zsh
./testappcmd
```