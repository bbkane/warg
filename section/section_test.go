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
			sec: section.New(
				"root section help",
				section.NewCommand("c1", "", nil, command.CompletionCandidates(nil)),
				section.NewSection("s1", "",
					section.NewCommand("c2", "", nil, command.CompletionCandidates(nil)),
				),
			),
			expected: []wargcore.FlatSection{
				{
					Path: []string{"r"},
					Sec: section.New(
						"root section help",
						section.NewCommand("c1", "", nil, command.CompletionCandidates(nil)),
						section.NewSection("s1", "",
							section.NewCommand("c2", "", nil, command.CompletionCandidates(nil)),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: section.New(
						"", section.NewCommand("c2", "", nil, command.CompletionCandidates(nil)),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "sortedOrder",
			rootName: "r",
			sec: section.New("",
				section.NewSection("sc", "",
					section.NewCommand("c", "", nil, command.CompletionCandidates(nil)),
				),
				section.NewSection("sb", "",
					section.NewCommand("c", "", nil, command.CompletionCandidates(nil)),
				),
				section.NewSection("sa", "",
					section.NewCommand("c", "", nil, command.CompletionCandidates(nil)),
				),
			),
			expected: []wargcore.FlatSection{
				{
					Path: []string{"r"},
					Sec: section.New("",
						section.NewSection("sc", "",
							section.NewCommand("c", "", nil, command.CompletionCandidates(nil)),
						),
						section.NewSection("sb", "",
							section.NewCommand("c", "", nil, command.CompletionCandidates(nil)),
						),
						section.NewSection("sa", "",
							section.NewCommand("c", "", nil, command.CompletionCandidates(nil)),
						),
					),
				},
				{
					Path: []string{"r", "sa"},
					Sec: section.New("",
						section.NewCommand("c", "", nil, command.CompletionCandidates(nil)),
					),
				},
				{
					Path: []string{"r", "sb"},
					Sec: section.New("",
						section.NewCommand("c", "", nil, command.CompletionCandidates(nil)),
					),
				},
				{
					Path: []string{"r", "sc"},
					Sec: section.New("",
						section.NewCommand("c", "", nil, command.CompletionCandidates(nil)),
					),
				},
			},
			expectedPanic: false,
		},
		{
			name:     "dupFlagNamesSeparatePaths",
			rootName: "r",
			sec: section.New("",
				section.NewSection("s1", "",
					section.NewCommand(
						"c1",
						"",
						nil,
						command.CompletionCandidates(nil),
						command.NewFlag("-f1", "", nil, flag.CompletionCandidates(nil)),
					),
				),
				section.NewSection("s2", "",
					section.NewCommand(
						"c1",
						"",
						nil,
						command.CompletionCandidates(nil),
						command.NewFlag("-f1", "", nil, flag.CompletionCandidates(nil)),
					),
				),
			),
			expected: []wargcore.FlatSection{
				{
					Path: []string{"r"},
					Sec: section.New("",
						section.NewSection("s1", "",
							section.NewCommand(
								"c1",
								"",
								nil,
								command.CompletionCandidates(nil),
								command.NewFlag("-f1", "", nil, flag.CompletionCandidates(nil)),
							),
						),
						section.NewSection("s2", "",
							section.NewCommand(
								"c1",
								"",
								nil,
								command.CompletionCandidates(nil),
								command.NewFlag("-f1", "", nil, flag.CompletionCandidates(nil)),
							),
						),
					),
				},
				{
					Path: []string{"r", "s1"},
					Sec: section.New("",
						section.NewCommand(
							"c1",
							"",
							nil,
							command.CompletionCandidates(nil),
							command.NewFlag("-f1", "", nil, flag.CompletionCandidates(nil)),
						),
					),
				},
				{
					Path: []string{"r", "s2"},
					Sec: section.New("",
						section.NewCommand(
							"c1",
							"",
							nil,
							command.CompletionCandidates(nil),
							command.NewFlag("-f1", "", nil, flag.CompletionCandidates(nil)),
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
