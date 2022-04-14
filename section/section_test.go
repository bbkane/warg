package section_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"go.bbkane.com/warg/flag"
	"go.bbkane.com/warg/section"
)

func TestSectionT_BreadthFirst(t *testing.T) {
	// NOTE: function equality cannot be compared with assert.Equal,
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
					Path:           []section.Name{"r"},
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
					Path:           []section.Name{"r", "s1"},
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
					Path:           []section.Name{"r"},
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
					Path:           []section.Name{"r", "sa"},
					Sec: section.New("",
						section.Command("c", "", nil),
					),
				},
				{
					InheritedFlags: make(flag.FlagMap),
					Path:           []section.Name{"r", "sb"},
					Sec: section.New("",
						section.Command("c", "", nil),
					),
				},
				{
					InheritedFlags: make(flag.FlagMap),
					Path:           []section.Name{"r", "sc"},
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
					Path:           []section.Name{"r"},
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
					Path:           []section.Name{"r", "s1"},
					Sec: section.New("",
						section.Command("c1", "", nil),
						section.Flag("-f1", "", nil),
					),
				},
				{
					InheritedFlags: make(flag.FlagMap),
					Path:           []section.Name{"r", "s2"},
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
				assert.Panics(
					t,
					func() {
						it := tt.sec.BreadthFirst([]section.Name{tt.rootName})
						for it.HasNext() {
							it.Next()
						}
					},
				)
				return
			}

			actual := make([]section.FlatSection, 0, 1)
			it := tt.sec.BreadthFirst([]section.Name{tt.rootName})
			for it.HasNext() {
				actual = append(actual, it.Next())
			}
			assert.Equal(t, tt.expected, actual)
		})
	}
}
