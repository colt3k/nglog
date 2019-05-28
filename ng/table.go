package ng

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	borderHorizCharacter = "-"
	borderVertCharacter  = "|"
)

/*
Add ability to pass in data of different types and convert to string in order to print
Add ability to pass in a json or csv file
*/

type TableOption func(h *PrinterTable)

func TableName(tblName string) TableOption {
	return func(pt *PrinterTable) {
		pt.tableName = tblName
	}
}
func TableHorizontal(horizontal bool) TableOption {
	return func(pt *PrinterTable) {
		pt.horizontal = horizontal
	}
}
func TableBorder(border bool) TableOption {
	return func(pt *PrinterTable) {
		pt.border = border
	}
}

/*
Add ability to pass in data of different types and convert to string in order to print
Add ability to pass in a json or csv file
*/

type Column struct {
	value string
}
type Row struct {
	colums []*Column
}
type Table struct {
	tableName string
	header    *Row
	rows      []*Row
}

type PrinterTable struct {
	Table
	horizontal         bool
	fullLength         int
	fullLengthPlusCols int
	szAr               []int
	border             bool
}

func NewTable(header []string, rows [][]string, opts ...TableOption) (*PrinterTable, error) {
	if len(header) != len(rows[0]) {
		return nil, errors.New("header and row columns should match, header: " + strconv.Itoa(len(header)) + ", row: " + strconv.Itoa(len(rows[0])))
	}

	t := new(PrinterTable)
	for _, opt := range opts {
		opt(t)
	}

	if !t.horizontal {
		t.populateTable(header, rows)
	} else {
		t.populateHorizontalColumn(header, rows)
	}
	return t, nil
}

func BuildRow(row []string) *Row {
	r := new(Row)
	cols := make([]*Column, 0)
	for _, d := range row {
		c := new(Column)
		c.value = " " + d + " "
		cols = append(cols, c)
	}
	r.colums = cols
	return r
}

func BuildTableBody(rowsData [][]string) []*Row {
	rows := make([]*Row, 0)
	for _, d := range rowsData {
		row := BuildRow(d)
		rows = append(rows, row)
	}

	return rows
}

//Header is first row and Rows are below
func (t *PrinterTable) populateTable(header []string, rows [][]string) {
	t.header = BuildRow(header)
	t.rows = BuildTableBody(rows)
}

//Header is first Column and Rows go across instead
func (t *PrinterTable) populateHorizontalColumn(header []string, rows [][]string) {

	rowData := make([]*Row, 0)
	//Re-organize data into columns instead of rows
	for _, d := range header {
		cols := make([]*Column, 0)
		c := new(Column)
		c.value = " " + d + " "
		cols = append(cols, c)

		r := new(Row)
		r.colums = cols
		rowData = append(rowData, r)
	}

	for _, d := range rows {
		for j, k := range d {
			//if i == 0 {
			r := rowData[j]
			cols := make([]*Column, 0)
			for _, t := range r.colums {
				cols = append(cols, t)
			}
			c := new(Column)
			c.value = " " + k + " "
			cols = append(cols, c)
			r.colums = cols
		}
	}

	t.rows = rowData
}

func (t *PrinterTable) setupBorder() {
	var w = 0

	if !t.horizontal {
		for _, d := range t.header.colums {
			rns := []rune(d.value)
			w = w + len(rns)
		}
	}
	for _, d := range t.rows {
		var tw = 0
		for _, d2 := range d.colums {
			rns := []rune(d2.value)
			tw = tw + len(rns)
		}
		if tw > w {
			w = tw
		}
	}
	t.fullLength = w
}

// Loop over each column and find largest width to make that column for all in it
func (t *PrinterTable) Print() {

	t.Sizes()
	t.AdjustData()
	fmt.Println()
	var borderVal string
	if t.border {
		t.setupBorder()
		cols := len(t.szAr) + 1
		t.fullLengthPlusCols = t.fullLength + cols
		borderVal = strings.Repeat(borderHorizCharacter, t.fullLengthPlusCols)
		fmt.Println(borderVal)
	}
	if len(t.tableName) > 0 {
		// split border length in half
		// split table name in half
		// name "Example" length = 7 / 2 = 3.5 first part end of 1/2 of border length
		//-----------------------------------|-----------------------------------
		//                                Example
		borderHalf := t.fullLengthPlusCols / 2
		rnsTblName := []rune(t.tableName)
		titleHalfSz := len(rnsTblName) / 2.0
		titleFirstHalf := rnsTblName[0:titleHalfSz]
		titleSecHalf := rnsTblName[titleHalfSz:]
		firstHalf := strings.Repeat(" ", borderHalf-titleHalfSz)
		accumP1 := (borderHalf - titleHalfSz) + len(titleSecHalf)
		wdthDiff := t.fullLengthPlusCols - accumP1
		//fmt.Println("AccumP1:", accumP1)
		//fmt.Println("Diff:", wdthDiff)
		//fmt.Printf("WdthTTL: %d calcTTL: %d\n", t.fullLengthPlusCols, (accumP1 + wdthDiff))
		var difOff int
		if accumP1 > wdthDiff {
			difOff = 1
		} else {
			difOff = 2
		}
		subSecHalf := wdthDiff - len(titleSecHalf) - difOff

		secondHalf := strings.Repeat(" ", subSecHalf)
		fmt.Print("|" + firstHalf + string(titleFirstHalf))
		fmt.Print(string(titleSecHalf) + secondHalf + "|")
		fmt.Println()

		//fmt.Print("|" + strings.Repeat(" ", borderHalf) + "|")
		//fmt.Println()
		fmt.Println(borderVal)
	}
	if !t.horizontal {
		for _, d := range t.header.colums {
			if t.border {
				fmt.Print(borderVertCharacter)
			}
			fmt.Print(d.value)
		}
		if t.border {
			fmt.Print(borderVertCharacter)
		}
		if t.border {
			fmt.Println()
			fmt.Println(borderVal)
		}
	}

	for _, d := range t.rows {
		for _, d2 := range d.colums {
			if t.border {
				fmt.Print(borderVertCharacter)
			}
			fmt.Print(d2.value)
		}
		if t.border {
			fmt.Print(borderVertCharacter)
		}
		if t.border {
			fmt.Println()
			fmt.Println(borderVal)
		} else {
			fmt.Println()
		}
	}

	fmt.Println()
}

func (t *PrinterTable) AdjustData() {
	if t.header != nil && len(t.header.colums) > 0 {
		for i, d := range t.header.colums {
			rns := []rune(d.value)
			sztmp := len(rns)
			sz := t.szAr[i]
			if sztmp < sz {
				//pad runes
				diff := sz - sztmp
				pad := strings.Repeat(" ", diff)
				c := new(Column)
				c.value = d.value + pad
				t.header.colums[i] = c
			}
		}
	}
	for i, d := range t.rows {
		for j, d2 := range d.colums {
			rns := []rune(d2.value)
			sztmp := len(rns)
			sz := t.szAr[j]
			if sztmp < sz {
				//pad runes
				diff := sz - sztmp
				pad := strings.Repeat(" ", diff)
				t.rows[i].colums[j].value = d2.value + pad
			}
		}
	}
}
func (t *PrinterTable) Sizes() {
	if t.header != nil && len(t.header.colums) > 0 {
		t.szAr = make([]int, len(t.header.colums))
	} else {
		t.szAr = make([]int, len(t.rows[0].colums))
	}
	if t.header != nil && len(t.header.colums) > 0 {
		for i, d := range t.header.colums {
			rns := []rune(d.value)
			sztmp := len(rns)
			if sztmp > t.szAr[i] {
				t.szAr[i] = sztmp
			}
		}
	}
	for _, d := range t.rows {
		for j, d2 := range d.colums {
			rns := []rune(d2.value)
			sztmp := len(rns)
			if sztmp > t.szAr[j] {
				t.szAr[j] = sztmp
			}
		}
	}
}
