package main

import (
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
)

func app() *warg.App {
	linkUnlinkFlags := flag.FlagMap{
		"--ask": flag.NewFlag(
			"Whether to ask before making changes",
			// value.StringEnum("true", "false", "dry-run"),
			scalar.String(
				scalar.Choices("true", "false", "dry-run"),
				scalar.Default("true"),
			),
			flag.Required(),
		),
		"--dotfiles": flag.NewFlag(
			"Files/dirs starting with 'dot-' will have links starting with '.'",
			scalar.Bool(
				scalar.Default(true),
			),
			flag.Required(),
		),
		"--ignore": flag.NewFlag(
			"Ignore file/dir if the name (not the whole path) matches passed regex",
			slice.String(
				slice.Default([]string{"README.*"}),
			),
			flag.Alias("-i"),
			flag.UnsetSentinel("UNSET"),
		),
		"--link-dir": flag.NewFlag(
			"Symlinks will be created in this directory pointing to files/directories in --src-dir",
			scalar.Path(
				scalar.Default(path.New("~")),
			),
			flag.Alias("-l"),
			flag.Required(),
		),
		"--src-dir": flag.NewFlag(
			"Directory containing files and directories to link to",
			scalar.Path(),
			flag.Alias("-s"),
			flag.Required(),
		),
	}

	app := warg.New(
		"fling",
		"v1.0.0",
		section.New(
			"Link and unlink directory heirarchies ",
			section.NewCommand(
				"link",
				"Create links",
				link,
				command.FlagMap(linkUnlinkFlags),
			),
			section.NewCommand(
				"unlink",
				"Unlink previously created links",
				unlink,
				command.FlagMap(linkUnlinkFlags),
			),
			section.CommandMap(warg.VersionCommandMap()),
			section.Footer("Homepage: https://github.com/bbkane/fling"),
		),
		warg.GlobalFlagMap(warg.ColorFlagMap()),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun()
}
