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
				warg.NewSubCmd("c1", "", nil),
				warg.NewSubSection("s1", "",
					warg.NewSubCmd("c2", "", nil),
				),
			),
			expected: []warg.FlatSection{
				{
					Path: []string{"r"},
					Sec: warg.NewSection(
						"root section help",
						warg.NewSubCmd("c1", "", nil),
						warg.NewSubSection("s1", "",
							warg.NewSubCmd("c2", "", nil),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: warg.NewSection(
						"", warg.NewSubCmd("c2", "", nil),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "sortedOrder",
			rootName: "r",
			sec: warg.NewSection("",
				warg.NewSubSection("sc", "",
					warg.NewSubCmd("c", "", nil),
				),
				warg.NewSubSection("sb", "",
					warg.NewSubCmd("c", "", nil),
				),
				warg.NewSubSection("sa", "",
					warg.NewSubCmd("c", "", nil),
				),
			),
			expected: []warg.FlatSection{
				{
					Path: []string{"r"},
					Sec: warg.NewSection("",
						warg.NewSubSection("sc", "",
							warg.NewSubCmd("c", "", nil),
						),
						warg.NewSubSection("sb", "",
							warg.NewSubCmd("c", "", nil),
						),
						warg.NewSubSection("sa", "",
							warg.NewSubCmd("c", "", nil),
						),
					),
				},
				{
					Path: []string{"r", "sa"},
					Sec: warg.NewSection("",
						warg.NewSubCmd("c", "", nil),
					),
				},
				{
					Path: []string{"r", "sb"},
					Sec: warg.NewSection("",
						warg.NewSubCmd("c", "", nil),
					),
				},
				{
					Path: []string{"r", "sc"},
					Sec: warg.NewSection("",
						warg.NewSubCmd("c", "", nil),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "dupFlagNamesSeparatePaths",
			rootName: "r",
			sec: warg.NewSection("",
				warg.NewSubSection("s1", "",
					warg.NewSubCmd(
						"c1",
						"",
						nil,
						warg.NewCmdFlag("-f1", "", nil, warg.FlagCompletions(nil)),
					),
				),
				warg.NewSubSection("s2", "",
					warg.NewSubCmd(
						"c1",
						"",
						nil,
						warg.NewCmdFlag("-f1", "", nil, warg.FlagCompletions(nil)),
					),
				),
			),
			expected: []warg.FlatSection{
				{
					Path: []string{"r"},
					Sec: warg.NewSection("",
						warg.NewSubSection("s1", "",
							warg.NewSubCmd(
								"c1",
								"",
								nil,
								warg.NewCmdFlag("-f1", "", nil, warg.FlagCompletions(nil)),
							),
						),
						warg.NewSubSection("s2", "",
							warg.NewSubCmd(
								"c1",
								"",
								nil,
								warg.NewCmdFlag("-f1", "", nil, warg.FlagCompletions(nil)),
							),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: warg.NewSection("",
						warg.NewSubCmd(
							"c1",
							"",
							nil,
							warg.NewCmdFlag("-f1", "", nil, warg.FlagCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "s2"},
					Sec: warg.NewSection("",
						warg.NewSubCmd(
							"c1",
							"",
							nil,
							warg.NewCmdFlag("-f1", "", nil, warg.FlagCompletions(nil)),
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
