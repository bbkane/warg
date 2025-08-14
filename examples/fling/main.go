package main

import (
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
	"go.bbkane.com/warg/wargcore"
)

func app() *wargcore.App {
	linkUnlinkFlags := wargcore.FlagMap{
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
		section.NewSection(
			"Link and unlink directory heirarchies ",
			section.NewChildCmd(
				"link",
				"Create links",
				link,
				command.ChildFlagMap(linkUnlinkFlags),
			),
			section.NewChildCmd(
				"unlink",
				"Unlink previously created links",
				unlink,
				command.ChildFlagMap(linkUnlinkFlags),
			),
			section.SectionFooter("Homepage: https://github.com/bbkane/fling"),
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun()
}
