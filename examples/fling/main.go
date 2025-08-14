package main

import (
	"go.bbkane.com/warg"
	"go.bbkane.com/warg/path"
	"go.bbkane.com/warg/value/scalar"
	"go.bbkane.com/warg/value/slice"
	"go.bbkane.com/warg/wargcore"
)

func app() *wargcore.App {
	linkUnlinkFlags := wargcore.FlagMap{
		"--ask": wargcore.NewFlag(
			"Whether to ask before making changes",
			// value.StringEnum("true", "false", "dry-run"),
			scalar.String(
				scalar.Choices("true", "false", "dry-run"),
				scalar.Default("true"),
			),
			wargcore.Required(),
		),
		"--dotfiles": wargcore.NewFlag(
			"Files/dirs starting with 'dot-' will have links starting with '.'",
			scalar.Bool(
				scalar.Default(true),
			),
			wargcore.Required(),
		),
		"--ignore": wargcore.NewFlag(
			"Ignore file/dir if the name (not the whole path) matches passed regex",
			slice.String(
				slice.Default([]string{"README.*"}),
			),
			wargcore.Alias("-i"),
			wargcore.UnsetSentinel("UNSET"),
		),
		"--link-dir": wargcore.NewFlag(
			"Symlinks will be created in this directory pointing to files/directories in --src-dir",
			scalar.Path(
				scalar.Default(path.New("~")),
			),
			wargcore.Alias("-l"),
			wargcore.Required(),
		),
		"--src-dir": wargcore.NewFlag(
			"Directory containing files and directories to link to",
			scalar.Path(),
			wargcore.Alias("-s"),
			wargcore.Required(),
		),
	}

	app := warg.New(
		"fling",
		"v1.0.0",
		wargcore.NewSection(
			"Link and unlink directory heirarchies ",
			wargcore.NewChildCmd(
				"link",
				"Create links",
				link,
				wargcore.ChildFlagMap(linkUnlinkFlags),
			),
			wargcore.NewChildCmd(
				"unlink",
				"Unlink previously created links",
				unlink,
				wargcore.ChildFlagMap(linkUnlinkFlags),
			),
			wargcore.SectionFooter("Homepage: https://github.com/bbkane/fling"),
		),
		warg.SkipValidation(),
	)
	return &app
}

func main() {
	app().MustRun()
}
