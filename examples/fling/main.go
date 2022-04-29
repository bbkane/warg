package main

import (
	"os"

	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value"
)

func app() *warg.App {
	linkUnlinkFlags := flag.FlagMap{
		"--ask": flag.New(
			"Whether to ask before making changes",
			value.StringEnum("true", "false", "dry-run"),
			flag.Default("true"),
			flag.Required(),
		),
		"--dotfiles": flag.New(
			"Files/dirs starting with 'dot-' will have links starting with '.'",
			value.Bool,
			flag.Default("true"),
			flag.Required(),
		),
		"--ignore": flag.New(
			"Ignore file/dir if the name (not the whole path) matches passed regex",
			value.StringSlice,
			flag.Alias("-i"),
		),
		"--link-dir": flag.New(
			"Symlinks will be created in this directory pointing to files/directories in --src-dir",
			value.Path,
			flag.Alias("-l"),
			flag.Default("~"),
			flag.Required(),
		),
		"--src-dir": flag.New(
			"Directory containing files and directories to link to",
			value.Path,
			flag.Alias("-s"),
			flag.Required(),
		),
	}

	app := warg.New(
		"fling",
		section.New(
			"Link and unlink directory heirarchies ",
			section.Command(
				"version",
				"Print version",
				printVersion,
			),
			section.Command(
				"link",
				"Create links",
				link,
				command.ExistingFlags(linkUnlinkFlags),
			),
			section.Command(
				"unlink",
				"Unlink previously created links",
				unlink,
				command.ExistingFlags(linkUnlinkFlags),
			),
			section.Flag(
				"--color",
				"Control color (including for --help)",
				value.StringEnum("true", "false", "auto"),
				flag.Alias("-c"),
				flag.Default("auto"),
			),
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun(os.Args, os.LookupEnv)
}
