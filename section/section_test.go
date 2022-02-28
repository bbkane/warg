package section_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
)

func TestSectionT_BreadthFirst(t *testing.T) {
	// NOTE: function equality cannot be compared with require.Equal,
	// so just set action to nil
	tests := []struct {
		name          string
		rootName      section.Name
		sec           section.SectionT
		expected      []section.FlatSection
		expectedPanic bool
	}{
		{
			name:     "simple",
			rootName: "r",
			sec: section.New(
				"root section help",
				section.Command("c1", "", nil),
				section.Section("s1", "",
					section.Command("c2", "", nil),
				),
			),
			expected: []section.FlatSection{
				{
					InheritedFlags: make(flag.FlagMap),
					Name:           "r",
					ParentPath:     []section.Name{},
					Sec: section.New(
						"root section help",
						section.Command("c1", "", nil),
						section.Section("s1", "",
							section.Command("c2", "", nil),
						),
					),
				},
				{
					InheritedFlags: make(flag.FlagMap),
					Name:           "s1",
					ParentPath:     []section.Name{"r"},
					Sec: section.New(
						"", section.Command("c2", "", nil),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "sortedOrder",
			rootName: "r",
			sec: section.New("",
				section.Section("sc", "",
					section.Command("c", "", nil),
				),
				section.Section("sb", "",
					section.Command("c", "", nil),
				),
				section.Section("sa", "",
					section.Command("c", "", nil),
				),
			),
			expected: []section.FlatSection{
				{
					InheritedFlags: make(flag.FlagMap),
					Name:           "r",
					ParentPath:     []section.Name{},
					Sec: section.New("",
						section.Section("sc", "",
							section.Command("c", "", nil),
						),
						section.Section("sb", "",
							section.Command("c", "", nil),
						),
						section.Section("sa", "",
							section.Command("c", "", nil),
						),
					),
				},
				{
					InheritedFlags: make(flag.FlagMap),
					Name:           "sa",
					ParentPath:     []section.Name{"r"},
					Sec: section.New("",
						section.Command("c", "", nil),
					),
				},
				{
					InheritedFlags: make(flag.FlagMap),
					Name:           "sb",
					ParentPath:     []section.Name{"r"},
					Sec: section.New("",
						section.Command("c", "", nil),
					),
				},
				{
					InheritedFlags: make(flag.FlagMap),
					Name:           "sc",
					ParentPath:     []section.Name{"r"},
					Sec: section.New("",
						section.Command("c", "", nil),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "duplicateFlagNames",
			rootName: "r",
			sec: section.New("root section help",
				section.Flag("-f1", "", nil),
				section.Section("s1", "",
					section.Flag("-f1", "", nil),
					section.Section("s2", ""),
				),
			),
			expected:      nil,
			expectedPanic: true,
		},
		{
			name:     "dupFlagNamesSeparatePaths",
			rootName: "r",
			sec: section.New("",
				section.Section("s1", "",
					section.Command("c1", "", nil),
					section.Flag("-f1", "", nil),
				),
				section.Section("s2", "",
					section.Command("c1", "", nil),
					section.Flag("-f1", "", nil),
				),
			),
			expected: []section.FlatSection{
				{
					InheritedFlags: make(flag.FlagMap),
					Name:           "r",
					ParentPath:     []section.Name{},
					Sec: section.New("",
						section.Section("s1", "",
							section.Command("c1", "", nil),
							section.Flag("-f1", "", nil),
						),
						section.Section("s2", "",
							section.Command("c1", "", nil),
							section.Flag("-f1", "", nil),
						),
					),
				},
				{
					InheritedFlags: make(flag.FlagMap),
					Name:           "s1",
					ParentPath:     []section.Name{"r"},
					Sec: section.New("",
						section.Command("c1", "", nil),
						section.Flag("-f1", "", nil),
					),
				},
				{
					InheritedFlags: make(flag.FlagMap),
					Name:           "s2",
					ParentPath:     []section.Name{"r"},
					Sec: section.New("",
						section.Command("c1", "", nil),
						section.Flag("-f1", "", nil),
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
						it := tt.sec.BreadthFirst(tt.rootName)
						for it.HasNext() {
							it.Next()
						}
					},
				)
				return
			}

			actual := make([]section.FlatSection, 0, 1)
			it := tt.sec.BreadthFirst(tt.rootName)
			for it.HasNext() {
				actual = append(actual, it.Next())
			}
			require.Equal(t, tt.expected, actual)
		})
	}
}
