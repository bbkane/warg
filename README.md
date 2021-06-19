See ~/journal/arg_parsing.md and ~/Git/bakeoff_argparse

# Naming Conventions

- `NewXXX`: Create an `XXX` from options
- `AddXXX`: Add a created `XXX` to your tree
- `WithXXX`: Create and add and `XXX` to your tree

# Next Steps

- function actions - see below
- --help
- Web form
- Config parsing

# Other library names

- argyoumeant
- chronicli
- clinch
- embargo
- gargle
- miracli
- monocli
- motorcycli
- oracli
- periclis
- receptacli
- target
- targpit
- tentacli
- warg
- yaargparse
- clide

# Other category names

- category
- namespace
- module
- commandgroup (az calls them groups)
- noun (with command -> verb rename too)

# Links

- https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
- https://dave.cheney.net/2016/11/13/do-not-fear-first-class-functions
- https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html

# Flags

I'm not sure the best way to do flags...

Kingpin uses https://pkg.go.dev/flag?utm_source=godoc#Getter - see https://github.com/alecthomas/kingpin#custom-parsers
urfave/cli uses different Structs for each flagtype and then shoves them into a context: https://github.com/urfave/cli/blob/master/docs/v2/manual.md#flags (NOTE: they also support getting flags from envvars, files, JSON like I want)

What about this:

```go
type Flag struct {
    Help string
    Parser func(interface{}, string) (error)
    Value interface{}
}
```

That's working for my IntListParser I think. I need to get that table driven and shove it into the flags. I also need to

Look out for https://medium.com/hackernoon/today-i-learned-pass-by-reference-on-interface-parameter-in-golang-35ee8d8a848e

So.. that did not work. See ~/Code/Go/interfaces for what I'm copying into the repo

# Parsing

from ~/journal/arg_parsing.md

- parse into a struct somehow (similar to JSON)
- parse level by level (similar to JSON)
- parse a whole line at once - similar to kingpin? `parsed.EnteredArgs("certificate create").Flag("id")`
- I wish I could generate an anonymous algebraic type and let people use that.

NOTE: I *can* make the parser create a new data structure for the user to use, instead of filling in the inline flags. can also separate FlagValueParams and PassedFlags into separate data structures

I'm not a huge fan of the following because it's easy to misspell a case statement and you probably don't know if you have all of them

```
parser := NewApp(...)
passedCommand, passedFlags, err := parser.RootCommand.Parse(os.Args[1:])
panicIf(err)
switch passedCommand{
case "cat1 com1":
    ...
default:
    ...
}
```

The following approach means I need to return a the LeafCommand of the parse tree and it needs to have a Execute function that can only return an error. So I don't like that

```
parser := NewApp(...)
leafCommand, passedFlags, err := parser.Parse(os.Args[1:])
panicIf(err)
err := passedCommand.Execute(passedFlags)
panicIf(err)
```

TODO: rename RootCommand to RootCategory, get better names for tests args (got and got1)

The perfect solution would let me loop through the possible LeafCommands, and manually call the functions I need with the passed FlagMap.
What if I had a method to loop through the possible commmands and then catch any that didn't get visited. We could put that in a test


Something like the following where GetCommand() sets a "visited value"

```
parser = NewApp(...)

switch leafCommand{
    case parser.GetCommand("cat com1"):
        ...
    case parser.GetCommand("cat1 com1"):
        ...
}

for command in parser.AllCommands():
    if not command.Visited():
        error
```

The action thing is easier so I'mma run with that for now :)
