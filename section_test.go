package warg_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"go.bbkane.com/warg"
)

func TestSectionT_BreadthFirst(t *testing.T) {
	// NOTE: function equality cannot be compared with assert.Equal,
	// so just set action to nil
	tests := []struct {
		name          string
		rootName      string
		sec           warg.Section
		expected      []warg.FlatSection
		expectedPanic bool
	}{
		{
			name:     "simple",
			rootName: "r",
			sec: warg.NewSection(
				"root section help",
				warg.NewChildCmd("c1", "", nil, warg.CmdCompletions(nil)),
				warg.NewChildSection("s1", "",
					warg.NewChildCmd("c2", "", nil, warg.CmdCompletions(nil)),
				),
			),
			expected: []warg.FlatSection{
				{
					Path: []string{"r"},
					Sec: warg.NewSection(
						"root section help",
						warg.NewChildCmd("c1", "", nil, warg.CmdCompletions(nil)),
						warg.NewChildSection("s1", "",
							warg.NewChildCmd("c2", "", nil, warg.CmdCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: warg.NewSection(
						"", warg.NewChildCmd("c2", "", nil, warg.CmdCompletions(nil)),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "sortedOrder",
			rootName: "r",
			sec: warg.NewSection("",
				warg.NewChildSection("sc", "",
					warg.NewChildCmd("c", "", nil, warg.CmdCompletions(nil)),
				),
				warg.NewChildSection("sb", "",
					warg.NewChildCmd("c", "", nil, warg.CmdCompletions(nil)),
				),
				warg.NewChildSection("sa", "",
					warg.NewChildCmd("c", "", nil, warg.CmdCompletions(nil)),
				),
			),
			expected: []warg.FlatSection{
				{
					Path: []string{"r"},
					Sec: warg.NewSection("",
						warg.NewChildSection("sc", "",
							warg.NewChildCmd("c", "", nil, warg.CmdCompletions(nil)),
						),
						warg.NewChildSection("sb", "",
							warg.NewChildCmd("c", "", nil, warg.CmdCompletions(nil)),
						),
						warg.NewChildSection("sa", "",
							warg.NewChildCmd("c", "", nil, warg.CmdCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "sa"},
					Sec: warg.NewSection("",
						warg.NewChildCmd("c", "", nil, warg.CmdCompletions(nil)),
					),
				},
				{
					Path: []string{"r", "sb"},
					Sec: warg.NewSection("",
						warg.NewChildCmd("c", "", nil, warg.CmdCompletions(nil)),
					),
				},
				{
					Path: []string{"r", "sc"},
					Sec: warg.NewSection("",
						warg.NewChildCmd("c", "", nil, warg.CmdCompletions(nil)),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "dupFlagNamesSeparatePaths",
			rootName: "r",
			sec: warg.NewSection("",
				warg.NewChildSection("s1", "",
					warg.NewChildCmd(
						"c1",
						"",
						nil,
						warg.CmdCompletions(nil),
						warg.NewChildFlag("-f1", "", nil, warg.FlagCompletions(nil)),
					),
				),
				warg.NewChildSection("s2", "",
					warg.NewChildCmd(
						"c1",
						"",
						nil,
						warg.CmdCompletions(nil),
						warg.NewChildFlag("-f1", "", nil, warg.FlagCompletions(nil)),
					),
				),
			),
			expected: []warg.FlatSection{
				{
					Path: []string{"r"},
					Sec: warg.NewSection("",
						warg.NewChildSection("s1", "",
							warg.NewChildCmd(
								"c1",
								"",
								nil,
								warg.CmdCompletions(nil),
								warg.NewChildFlag("-f1", "", nil, warg.FlagCompletions(nil)),
							),
						),
						warg.NewChildSection("s2", "",
							warg.NewChildCmd(
								"c1",
								"",
								nil,
								warg.CmdCompletions(nil),
								warg.NewChildFlag("-f1", "", nil, warg.FlagCompletions(nil)),
							),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: warg.NewSection("",
						warg.NewChildCmd(
							"c1",
							"",
							nil,
							warg.CmdCompletions(nil),
							warg.NewChildFlag("-f1", "", nil, warg.FlagCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "s2"},
					Sec: warg.NewSection("",
						warg.NewChildCmd(
							"c1",
							"",
							nil,
							warg.CmdCompletions(nil),
							warg.NewChildFlag("-f1", "", nil, warg.FlagCompletions(nil)),
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

			actual := make([]warg.FlatSection, 0, 1)
			it := tt.sec.BreadthFirst([]string{tt.rootName})
			for it.HasNext() {
				actual = append(actual, it.Next())
			}
			require.Equal(t, tt.expected, actual)
		})
	}
}
