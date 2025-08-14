package wargcore_test

import (
	"testing"

	"github.com/stretchr/testify/require"

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
			sec: wargcore.NewSection(
				"root section help",
				wargcore.NewChildCmd("c1", "", nil, wargcore.CmdCompletions(nil)),
				wargcore.NewChildSection("s1", "",
					wargcore.NewChildCmd("c2", "", nil, wargcore.CmdCompletions(nil)),
				),
			),
			expected: []wargcore.FlatSection{
				{
					Path: []string{"r"},
					Sec: wargcore.NewSection(
						"root section help",
						wargcore.NewChildCmd("c1", "", nil, wargcore.CmdCompletions(nil)),
						wargcore.NewChildSection("s1", "",
							wargcore.NewChildCmd("c2", "", nil, wargcore.CmdCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: wargcore.NewSection(
						"", wargcore.NewChildCmd("c2", "", nil, wargcore.CmdCompletions(nil)),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "sortedOrder",
			rootName: "r",
			sec: wargcore.NewSection("",
				wargcore.NewChildSection("sc", "",
					wargcore.NewChildCmd("c", "", nil, wargcore.CmdCompletions(nil)),
				),
				wargcore.NewChildSection("sb", "",
					wargcore.NewChildCmd("c", "", nil, wargcore.CmdCompletions(nil)),
				),
				wargcore.NewChildSection("sa", "",
					wargcore.NewChildCmd("c", "", nil, wargcore.CmdCompletions(nil)),
				),
			),
			expected: []wargcore.FlatSection{
				{
					Path: []string{"r"},
					Sec: wargcore.NewSection("",
						wargcore.NewChildSection("sc", "",
							wargcore.NewChildCmd("c", "", nil, wargcore.CmdCompletions(nil)),
						),
						wargcore.NewChildSection("sb", "",
							wargcore.NewChildCmd("c", "", nil, wargcore.CmdCompletions(nil)),
						),
						wargcore.NewChildSection("sa", "",
							wargcore.NewChildCmd("c", "", nil, wargcore.CmdCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "sa"},
					Sec: wargcore.NewSection("",
						wargcore.NewChildCmd("c", "", nil, wargcore.CmdCompletions(nil)),
					),
				},
				{
					Path: []string{"r", "sb"},
					Sec: wargcore.NewSection("",
						wargcore.NewChildCmd("c", "", nil, wargcore.CmdCompletions(nil)),
					),
				},
				{
					Path: []string{"r", "sc"},
					Sec: wargcore.NewSection("",
						wargcore.NewChildCmd("c", "", nil, wargcore.CmdCompletions(nil)),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "dupFlagNamesSeparatePaths",
			rootName: "r",
			sec: wargcore.NewSection("",
				wargcore.NewChildSection("s1", "",
					wargcore.NewChildCmd(
						"c1",
						"",
						nil,
						wargcore.CmdCompletions(nil),
						wargcore.NewChildFlag("-f1", "", nil, wargcore.FlagCompletions(nil)),
					),
				),
				wargcore.NewChildSection("s2", "",
					wargcore.NewChildCmd(
						"c1",
						"",
						nil,
						wargcore.CmdCompletions(nil),
						wargcore.NewChildFlag("-f1", "", nil, wargcore.FlagCompletions(nil)),
					),
				),
			),
			expected: []wargcore.FlatSection{
				{
					Path: []string{"r"},
					Sec: wargcore.NewSection("",
						wargcore.NewChildSection("s1", "",
							wargcore.NewChildCmd(
								"c1",
								"",
								nil,
								wargcore.CmdCompletions(nil),
								wargcore.NewChildFlag("-f1", "", nil, wargcore.FlagCompletions(nil)),
							),
						),
						wargcore.NewChildSection("s2", "",
							wargcore.NewChildCmd(
								"c1",
								"",
								nil,
								wargcore.CmdCompletions(nil),
								wargcore.NewChildFlag("-f1", "", nil, wargcore.FlagCompletions(nil)),
							),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: wargcore.NewSection("",
						wargcore.NewChildCmd(
							"c1",
							"",
							nil,
							wargcore.CmdCompletions(nil),
							wargcore.NewChildFlag("-f1", "", nil, wargcore.FlagCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "s2"},
					Sec: wargcore.NewSection("",
						wargcore.NewChildCmd(
							"c1",
							"",
							nil,
							wargcore.CmdCompletions(nil),
							wargcore.NewChildFlag("-f1", "", nil, wargcore.FlagCompletions(nil)),
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
