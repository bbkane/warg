package main

import (
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
)

func app() *warg.App {
	linkUnlinkFlags := warg.FlagMap{
		"--ask": warg.NewFlag(
			"Whether to ask before making changes",
			// value.StringEnum("true", "false", "dry-run"),
			scalar.String(
				scalar.Choices("true", "false", "dry-run"),
				scalar.Default("true"),
			),
			warg.Required(),
		),
		"--dotfiles": warg.NewFlag(
			"Files/dirs starting with 'dot-' will have links starting with '.'",
			scalar.Bool(
				scalar.Default(true),
			),
			warg.Required(),
		),
		"--ignore": warg.NewFlag(
			"Ignore file/dir if the name (not the whole path) matches passed regex",
			slice.String(
				slice.Default([]string{"README.*"}),
			),
			warg.Alias("-i"),
			warg.UnsetSentinel("UNSET"),
		),
		"--link-dir": warg.NewFlag(
			"Symlinks will be created in this directory pointing to files/directories in --src-dir",
			scalar.Path(
				scalar.Default(path.New("~")),
			),
			warg.Alias("-l"),
			warg.Required(),
		),
		"--src-dir": warg.NewFlag(
			"Directory containing files and directories to link to",
			scalar.Path(),
			warg.Alias("-s"),
			warg.Required(),
		),
	}

	app := warg.New(
		"fling",
		"v1.0.0",
		warg.NewSection(
			"Link and unlink directory heirarchies ",
			warg.NewSubCmd(
				"link",
				"Create links",
				link,
				warg.CmdFlagMap(linkUnlinkFlags),
			),
			warg.NewSubCmd(
				"unlink",
				"Unlink previously created links",
				unlink,
				warg.CmdFlagMap(linkUnlinkFlags),
			),
			warg.SectionFooter("Homepage: https://github.com/bbkane/fling"),
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun()
}
