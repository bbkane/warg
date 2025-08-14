package section_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/command"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
	"go.bbkane.com/warg/wargcore"
)

func TestSectionT_BreadthFirst(t *testing.T) {
	// NOTE: function equality cannot be compared with assert.Equal,
	// so just set action to nil
	tests := []struct {
		name          string
		rootName      string
		sec           wargcore.Section
		expected      []wargcore.FlatSection
		expectedPanic bool
	}{
		{
			name:     "simple",
			rootName: "r",
			sec: section.NewSection(
				"root section help",
				section.NewChildCmd("c1", "", nil, command.CmdCompletions(nil)),
				section.NewChildSection("s1", "",
					section.NewChildCmd("c2", "", nil, command.CmdCompletions(nil)),
				),
			),
			expected: []wargcore.FlatSection{
				{
					Path: []string{"r"},
					Sec: section.NewSection(
						"root section help",
						section.NewChildCmd("c1", "", nil, command.CmdCompletions(nil)),
						section.NewChildSection("s1", "",
							section.NewChildCmd("c2", "", nil, command.CmdCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: section.NewSection(
						"", section.NewChildCmd("c2", "", nil, command.CmdCompletions(nil)),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "sortedOrder",
			rootName: "r",
			sec: section.NewSection("",
				section.NewChildSection("sc", "",
					section.NewChildCmd("c", "", nil, command.CmdCompletions(nil)),
				),
				section.NewChildSection("sb", "",
					section.NewChildCmd("c", "", nil, command.CmdCompletions(nil)),
				),
				section.NewChildSection("sa", "",
					section.NewChildCmd("c", "", nil, command.CmdCompletions(nil)),
				),
			),
			expected: []wargcore.FlatSection{
				{
					Path: []string{"r"},
					Sec: section.NewSection("",
						section.NewChildSection("sc", "",
							section.NewChildCmd("c", "", nil, command.CmdCompletions(nil)),
						),
						section.NewChildSection("sb", "",
							section.NewChildCmd("c", "", nil, command.CmdCompletions(nil)),
						),
						section.NewChildSection("sa", "",
							section.NewChildCmd("c", "", nil, command.CmdCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "sa"},
					Sec: section.NewSection("",
						section.NewChildCmd("c", "", nil, command.CmdCompletions(nil)),
					),
				},
				{
					Path: []string{"r", "sb"},
					Sec: section.NewSection("",
						section.NewChildCmd("c", "", nil, command.CmdCompletions(nil)),
					),
				},
				{
					Path: []string{"r", "sc"},
					Sec: section.NewSection("",
						section.NewChildCmd("c", "", nil, command.CmdCompletions(nil)),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "dupFlagNamesSeparatePaths",
			rootName: "r",
			sec: section.NewSection("",
				section.NewChildSection("s1", "",
					section.NewChildCmd(
						"c1",
						"",
						nil,
						command.CmdCompletions(nil),
						command.NewChildFlag("-f1", "", nil, flag.FlagCompletions(nil)),
					),
				),
				section.NewChildSection("s2", "",
					section.NewChildCmd(
						"c1",
						"",
						nil,
						command.CmdCompletions(nil),
						command.NewChildFlag("-f1", "", nil, flag.FlagCompletions(nil)),
					),
				),
			),
			expected: []wargcore.FlatSection{
				{
					Path: []string{"r"},
					Sec: section.NewSection("",
						section.NewChildSection("s1", "",
							section.NewChildCmd(
								"c1",
								"",
								nil,
								command.CmdCompletions(nil),
								command.NewChildFlag("-f1", "", nil, flag.FlagCompletions(nil)),
							),
						),
						section.NewChildSection("s2", "",
							section.NewChildCmd(
								"c1",
								"",
								nil,
								command.CmdCompletions(nil),
								command.NewChildFlag("-f1", "", nil, flag.FlagCompletions(nil)),
							),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: section.NewSection("",
						section.NewChildCmd(
							"c1",
							"",
							nil,
							command.CmdCompletions(nil),
							command.NewChildFlag("-f1", "", nil, flag.FlagCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "s2"},
					Sec: section.NewSection("",
						section.NewChildCmd(
							"c1",
							"",
							nil,
							command.CmdCompletions(nil),
							command.NewChildFlag("-f1", "", nil, flag.FlagCompletions(nil)),
						),
					),
				},
			},
			expectedPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.expectedPanic {
				require.Panics(
					t,
					func() {
						it := tt.sec.BreadthFirst([]string{tt.rootName})
						for it.HasNext() {
							it.Next()
						}
					},
				)
				return
			}

			actual := make([]wargcore.FlatSection, 0, 1)
			it := tt.sec.BreadthFirst([]string{tt.rootName})
			for it.HasNext() {
				actual = append(actual, it.Next())
			}
			require.Equal(t, tt.expected, actual)
		})
	}
}
