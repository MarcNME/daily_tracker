package table

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable_AddRow(t1 *testing.T) {
	tests := []struct {
		name   string
		table  Table
		want   Table
		newRow []string
	}{
		{
			name: "AddRow",
			table: Table{
				headRow:  []string{},
				rows:     [][]string{{"123", "456", "789"}},
				cellsLen: []int{3, 3, 3},
			},
			want: Table{
				headRow:  []string{},
				rows:     [][]string{{"123", "456", "789"}, {"abc", "def", "ghi"}},
				cellsLen: []int{3, 3, 3},
			},
			newRow: []string{"abc", "def", "ghi"},
		},
		{
			name: "AddRow longer cell",
			table: Table{
				headRow:  []string{},
				rows:     [][]string{{"123", "456", "789"}},
				cellsLen: []int{3, 3, 3},
			},
			want: Table{
				headRow:  []string{},
				rows:     [][]string{{"123", "456", "789"}, {"abc", "def", "ghij"}},
				cellsLen: []int{3, 3, 4},
			},
			newRow: []string{"abc", "def", "ghij"},
		},
		{
			name: "AddRow longer row",
			table: Table{
				headRow:  []string{},
				rows:     [][]string{{"123", "456", "789"}},
				cellsLen: []int{3, 3, 3},
			},
			want: Table{
				headRow:  []string{},
				rows:     [][]string{{"123", "456", "789"}, {"abc", "def", "ghi", "jkl"}},
				cellsLen: []int{3, 3, 3, 3},
			},
			newRow: []string{"abc", "def", "ghi", "jkl"},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			tt.table.AddRow(tt.newRow...)

			assert.Equal(t1, tt.want, tt.table)
		})
	}
}

func TestNewTable(t *testing.T) {
	table := NewTable()
	assert.NotNil(t, table.rows)
	assert.NotNil(t, table.cellsLen)
	assert.NotNil(t, table.headRow)
}

func TestTable_SetHeader(t *testing.T) {
	tests := []struct {
		name      string
		table     Table
		want      Table
		newHeader []string
	}{
		{
			name: "SetHeader",
			table: Table{
				headRow:  []string{},
				rows:     [][]string{{"123", "456", "789"}},
				cellsLen: []int{3, 3, 3},
			},
			want: Table{
				headRow:  []string{"abc", "def", "ghi"},
				rows:     [][]string{{"123", "456", "789"}},
				cellsLen: []int{3, 3, 3},
			},
			newHeader: []string{"abc", "def", "ghi"},
		},
		{
			name: "SetHeader replace old header",
			table: Table{
				headRow:  []string{"old", "header", "text"},
				rows:     [][]string{{"123", "456", "789"}},
				cellsLen: []int{3, 6, 4},
			},
			want: Table{
				headRow:  []string{"abc", "def", "ghi"},
				rows:     [][]string{{"123", "456", "789"}},
				cellsLen: []int{3, 6, 4},
			},
			newHeader: []string{"abc", "def", "ghi"},
		},
		{
			name: "SetHeader replace old header longer cells",
			table: Table{
				headRow:  []string{"a", "b", "c"},
				rows:     [][]string{{"123", "456", "789"}},
				cellsLen: []int{3, 3, 3},
			},
			want: Table{
				headRow:  []string{"new", "longer", "header", "text"},
				rows:     [][]string{{"123", "456", "789"}},
				cellsLen: []int{3, 6, 6, 4},
			},
			newHeader: []string{"new", "longer", "header", "text"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			tt.table.SetHeader(tt.newHeader...)

			assert.Equal(t1, tt.want, tt.table)
		})
	}
}

func TestTable_ToString(t *testing.T) {
	tests := []struct {
		name  string
		table Table
		want  string
	}{
		{
			name: "Only header",
			table: Table{
				headRow:  []string{"abc", "def", "ghi"},
				rows:     [][]string{},
				cellsLen: []int{3, 3, 3},
			},
			want: "abcdefghi\n",
		},
		{
			name: "Only header with column seperator",
			table: Table{
				headRow:         []string{"abc", "def", "ghi"},
				rows:            [][]string{},
				cellsLen:        []int{3, 3, 3},
				columnSeparator: " ",
			},
			want: "abc def ghi\n",
		},
		{
			name: "Only header with padding",
			table: Table{
				headRow:  []string{"abc", "def", "ghi"},
				rows:     [][]string{},
				cellsLen: []int{3, 3, 3},
				padding:  2,
			},
			want: "  abc    def    ghi\n",
		},
		{
			name: "All",
			table: Table{
				headRow:         []string{"abc", "def", "ghi"},
				rows:            [][]string{{"cell", "cell", "cell"}},
				cellsLen:        []int{4, 4, 4},
				columnSeparator: "|",
			},
			want: "abc |def |ghi\ncell|cell|cell\n",
		},
		{
			name: "header seperator different line length",
			table: Table{
				headRow:         []string{"abc", "def", "ghi"},
				rows:            [][]string{{"cell1", "cell2", "cell3"}},
				cellsLen:        []int{5, 5, 5},
				columnSeparator: "|",
				headerSeparator: "_",
			},
			want: "abc  |def  |ghi\n_________________\ncell1|cell2|cell3\n",
		},
		{
			name: "header seperator with padding",
			table: Table{
				headRow:         []string{"abc1", "def2", "ghi3"},
				rows:            [][]string{{"cell1", "cell2", "cell3"}},
				cellsLen:        []int{5, 5, 5},
				columnSeparator: "|",
				headerSeparator: "-",
				padding:         2,
			},
			want: "  abc1   |  def2   |  ghi3\n---------------------------\n  cell1  |  cell2  |  cell3\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t1 *testing.T) {
			result := tt.table.ToString()
			fmt.Println(result)
			assert.Equal(t1, tt.want, result)
		})
	}
}
