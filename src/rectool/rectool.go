package rectool

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	REC_EOL = "!"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

const (
	TEXT   string = "_"
	NUMBER string = "#"
)

type Column struct {
	name            string
	dataType        string
	questionColumn  int
	questionLine    int
	questionColor   int
	fieldColumn     int
	fieldLine       int
	fieldType       string
	fieldWidth      int
	entryFieldColor int
	description     string
	intValue        int
	strValue        string
}

type RECHeader struct {
	count   int
	length  int
	columns []Column
}

type Record struct {
	columns []Column
}

type RECFile struct {
	Header   RECHeader
	Records  []Record
	Filename string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func Load(filename string) {
	file, err := os.Open(filename)
	check(err)

	// TODO: support file encodings
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	file.Close()

	header := newHeader(text)
	body := newBody(header, text)

	//fmt.Printf(InfoColor, "********************************************************")
	//fmt.Println("")
	//fmt.Printf(InfoColor, "************************ Header ************************")
	//fmt.Println("")
	//fmt.Printf(InfoColor, "********************************************************")
	//fmt.Println("")
	//
	//fmt.Println(header)

	fmt.Printf(InfoColor, "********************************************************")
	fmt.Println("")
	fmt.Printf(InfoColor, "************************* Body *************************")
	fmt.Println("")
	fmt.Printf(InfoColor, "********************************************************")
	fmt.Println("")
	fmt.Println(body)

}

/*
 * Builds a new RECHeader based on the incoming
 * string array.
 */
func newHeader(rows []string) *RECHeader {
	row := strings.Fields(rows[0])
	if len(row) > 2 {
		panic("Malformed file: first line is wrong")
	}

	count, err := strconv.Atoi(row[0])
	check(err)

	h := RECHeader{count: count, length: 0}

	for i := 1; i <= h.count; i++ {
		row := rows[i]
		column := newHeaderColumn(row)
		h.columns = append(h.columns, column)
		h.length += column.fieldWidth
	}

	return &h
}

/*
 * Builds a new Column based on the incoming
 * string.
 */
func newHeaderColumn(row string) Column {
	words := strings.Fields(row)

	if len(words) < 10 {
		panic("Malformed REC header")
	}

	var hc = Column{}
	hc.name = words[0]
	hc.dataType = words[1]
	hc.questionColumn, _ = strconv.Atoi(words[2])
	hc.questionColor, _ = strconv.Atoi(words[3])
	hc.fieldColumn, _ = strconv.Atoi(words[4])
	hc.fieldLine, _ = strconv.Atoi(words[5])
	hc.fieldType = words[6]
	hc.fieldWidth, _ = strconv.Atoi(words[7])
	hc.entryFieldColor, _ = strconv.Atoi(words[8])
	hc.description = words[9]

	return hc
}

/*
 * Builds a new body based on the incoming string array
 */
func newBody(h *RECHeader, rows []string) *[]Record {
	var records []Record

	fmt.Printf(InfoColor, strconv.Itoa(h.count))

	dataIndex := h.count + 1
	data := strings.ReplaceAll(strings.Join(rows[dataIndex:], ""), REC_EOL, "")

	fmt.Print(DebugColor, "Length --> "+strconv.Itoa(h.length))
	fmt.Println("")

	rowCount := len(data) / h.length

	for rowIndex := 0; rowIndex < rowCount; rowIndex++ {

		var record = Record{}
		rowStartPosition := rowIndex * h.length
		rowEndPosition := rowIndex*h.length + h.length - 1
		chunk := data[rowStartPosition:rowEndPosition]

		fieldIndex := 0
		for i := 0; i < h.count; i++ {
			column := h.columns[i]

			// Skip not printed columns
			if column.fieldWidth == 0 {
				continue
			}

			fieldStartPosition := fieldIndex
			fieldEndPosition := fieldIndex + column.fieldWidth - 1
			fieldValue := chunk[fieldStartPosition:fieldEndPosition]

			if column.name[0:1] == TEXT {
				column.strValue = fieldValue
			} else {
				column.intValue, _ = strconv.Atoi(fieldValue)
			}

			fieldIndex += column.fieldWidth
			record.columns = append(record.columns, column)
		}

		records = append(records, record)
	}

	return &records
}
