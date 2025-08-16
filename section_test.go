package warg

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSectionT_BreadthFirst(t *testing.T) {
	// NOTE: function equality cannot be compared with assert.Equal,
	// so just set action to nil
	tests := []struct {
		name          string
		rootName      string
		sec           Section
		expected      []flatSection
		expectedPanic bool
	}{
		{
			name:     "simple",
			rootName: "r",
			sec: NewSection(
				"root section help",
				NewSubCmd("c1", "", nil),
				NewSubSection("s1", "",
					NewSubCmd("c2", "", nil),
				),
			),
			expected: []flatSection{
				{
					Path: []string{"r"},
					Sec: NewSection(
						"root section help",
						NewSubCmd("c1", "", nil),
						NewSubSection("s1", "",
							NewSubCmd("c2", "", nil),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: NewSection(
						"", NewSubCmd("c2", "", nil),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "sortedOrder",
			rootName: "r",
			sec: NewSection("",
				NewSubSection("sc", "",
					NewSubCmd("c", "", nil),
				),
				NewSubSection("sb", "",
					NewSubCmd("c", "", nil),
				),
				NewSubSection("sa", "",
					NewSubCmd("c", "", nil),
				),
			),
			expected: []flatSection{
				{
					Path: []string{"r"},
					Sec: NewSection("",
						NewSubSection("sc", "",
							NewSubCmd("c", "", nil),
						),
						NewSubSection("sb", "",
							NewSubCmd("c", "", nil),
						),
						NewSubSection("sa", "",
							NewSubCmd("c", "", nil),
						),
					),
				},
				{
					Path: []string{"r", "sa"},
					Sec: NewSection("",
						NewSubCmd("c", "", nil),
					),
				},
				{
					Path: []string{"r", "sb"},
					Sec: NewSection("",
						NewSubCmd("c", "", nil),
					),
				},
				{
					Path: []string{"r", "sc"},
					Sec: NewSection("",
						NewSubCmd("c", "", nil),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "dupFlagNamesSeparatePaths",
			rootName: "r",
			sec: NewSection("",
				NewSubSection("s1", "",
					NewSubCmd(
						"c1",
						"",
						nil,
						NewCmdFlag("-f1", "", nil, FlagCompletions(nil)),
					),
				),
				NewSubSection("s2", "",
					NewSubCmd(
						"c1",
						"",
						nil,
						NewCmdFlag("-f1", "", nil, FlagCompletions(nil)),
					),
				),
			),
			expected: []flatSection{
				{
					Path: []string{"r"},
					Sec: NewSection("",
						NewSubSection("s1", "",
							NewSubCmd(
								"c1",
								"",
								nil,
								NewCmdFlag("-f1", "", nil, FlagCompletions(nil)),
							),
						),
						NewSubSection("s2", "",
							NewSubCmd(
								"c1",
								"",
								nil,
								NewCmdFlag("-f1", "", nil, FlagCompletions(nil)),
							),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: NewSection("",
						NewSubCmd(
							"c1",
							"",
							nil,
							NewCmdFlag("-f1", "", nil, FlagCompletions(nil)),
						),
					),
				},
				{
					Path: []string{"r", "s2"},
					Sec: NewSection("",
						NewSubCmd(
							"c1",
							"",
							nil,
							NewCmdFlag("-f1", "", nil, FlagCompletions(nil)),
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
						it := tt.sec.breadthFirst([]string{tt.rootName})
						for it.HasNext() {
							it.Next()
						}
					},
				)
				return
			}

			actual := make([]flatSection, 0, 1)
			it := tt.sec.breadthFirst([]string{tt.rootName})
			for it.HasNext() {
				actual = append(actual, it.Next())
			}
			require.Equal(t, tt.expected, actual)
		})
	}
}
