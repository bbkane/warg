# DO NOT USE

This library is still under heavy development, will probably break, and will definitely change.

See ~/journal/arg_parsing.md and ~/Git/bakeoff_argparse

# Naming Conventions

- `NewXXX`: Create an `XXX` from options
- `AddXXX`: Add a created `XXX` to your tree
- `WithXXX`: Create and add an `XXX` to your tree

# TODO: Next milestone: grabbit

- make colored help
- make help less verbose...
- add flag.TypeInfo to be set to the value.TypeInfo - then I can make default a list when appropriate
- print long values more neatly.... make a value.StringList() method?
- figure out what to do with --color , --help requires it?
- Write good doc comments, README, examples
- make an app.Test() method folks can add to their apps - should test for unique flag names between parent and child sections/commands for one thing
- go through TODOs
- --help ideas: man, json, web, form, term, lsp, bash-completion, zsh-completion, outline
