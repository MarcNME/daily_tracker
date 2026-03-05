package table

import (
	"fmt"
	"strings"
)

type ITable interface {
	ToString() string
	AddRow(...string)
	DeleteRows()
	SetHeader(...string)
	SetPadding(int) error
	SetRowSeparator(string)
	SetColumnSeparator(string)
	SetHeaderSeparator(string)
}

type Table struct {
	headRow         []string
	rows            [][]string
	cellsLen        []int
	padding         int
	columnSeparator string
	rowSeparator    string
	headerSeparator string
}

func NewTable() Table {
	return Table{
		headRow:         make([]string, 0),
		rows:            make([][]string, 0),
		cellsLen:        make([]int, 0),
		columnSeparator: " ",
	}
}

func (t *Table) ToString() string {
	sb := &strings.Builder{}

	t.rowToString(sb, t.headRow)
	sb.WriteString("\n")

	if t.headerSeparator != "" {
		lineLength := 0
		for i, l := range t.cellsLen {
			lineLength += l
			if i != len(t.cellsLen)-1 {
				lineLength += t.padding * 2
				lineLength += len(t.columnSeparator)
			} else {
				lineLength += t.padding
			}

		}

		sb.WriteString(strings.Repeat(t.headerSeparator, lineLength))
		sb.WriteString("\n")
	}

	for _, row := range t.rows {
		t.rowToString(sb, row)
		sb.WriteString("\n")
	}

	return sb.String()
}

func (t *Table) rowToString(sb *strings.Builder, row []string) {
	for i, cell := range row {
		if t.padding != 0 {
			sb.WriteString(strings.Repeat(" ", t.padding))
		}
		sb.WriteString(cell)

		if i != len(row)-1 {
			if t.cellsLen[i] > len(cell) {
				sb.WriteString(strings.Repeat(" ", t.cellsLen[i]-len(cell)))
			}
			if t.padding != 0 {
				sb.WriteString(strings.Repeat(" ", t.padding))
			}
			sb.WriteString(t.columnSeparator)
		}
	}
}

func (t *Table) AddRow(row ...string) {
	t.rows = append(t.rows, row)

	t.adjustCellLength(row)
}

func (t *Table) SetHeader(header ...string) {
	t.headRow = header

	t.adjustCellLength(header)
}

func (t *Table) adjustCellLength(newRow []string) {
	if len(newRow) > len(t.cellsLen) {
		for i := 0; i <= len(newRow)-len(t.cellsLen); i++ {
			t.cellsLen = append(t.cellsLen, 0)
		}
	}
	for i, cell := range newRow {
		if t.cellsLen[i] < len(cell) {
			t.cellsLen[i] = len(cell)
		}
	}
}

func (t *Table) SetPadding(n int) error {
	if n < 0 {
		return fmt.Errorf("padding cannot be less than zero")
	}
	t.padding = n
	return nil
}

func (t *Table) SetRowSeparator(s string) {
	t.rowSeparator = s
}

func (t *Table) SetColumnSeparator(s string) {
	t.columnSeparator = s
}

func (t *Table) SetHeaderSeparator(s string) {
	t.headerSeparator = s
}

func (t *Table) DeleteRows() {
	length := len(t.rows)
	t.rows = make([][]string, 0, length)
}
