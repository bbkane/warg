# Testing

There are two "parts" of completion:

 - the selection parts - tested in `app_completion_ext_test.go`
 - the zsh integration parts - tested manually.

 ## `zsh` Integration Tests

ChatGPT offers ways to unit test zsh completion, and there's a few examples on StackOverflow, but I'll just manually test for now till I feel up to trying to automate it.

Setup:

```zsh
# install the app to $GOBIN
go install ./completion/internal/testappcmd
cd ~/go/bin  # default $GOBIN
# install the completion to something on $FPATH
./testappcmd --completion-script-zsh > ~/fbin/_testappcmd
# open a new shell to load completions
./testappcmd ...  # tab away!!
```

Run:

```zsh
./testappcmd manual TAB
```