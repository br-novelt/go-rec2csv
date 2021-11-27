package rectool

import (
	"bufio"
	"encoding/csv"
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
	value           string
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
		fmt.Printf(ErrorColor, "[ERROR] Something went wrong")
		fmt.Println("")
		panic(e)
	}
}

func Load(filename string) RECFile {
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

	content := RECFile{Header: header, Records: body, Filename: filename}

	fmt.Printf(InfoColor, "[COMPLETE] REC file loading")
	fmt.Println("")

	return content
}

/*
 * Builds a new RECHeader based on the incoming
 * string array.
 */
func newHeader(rows []string) RECHeader {
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

	return h
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
func newBody(h RECHeader, rows []string) []Record {
	var records []Record

	dataIndex := h.count + 1
	data := strings.ReplaceAll(strings.Join(rows[dataIndex:], ""), REC_EOL, "")
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

			// Last position reaches end of line
			var fieldEndPosition int
			if (fieldIndex + column.fieldWidth) >= h.length {
				fieldEndPosition = fieldIndex + column.fieldWidth - 1
			} else {
				fieldEndPosition = fieldIndex + column.fieldWidth
			}

			fieldStartPosition := fieldIndex
			fieldValue := chunk[fieldStartPosition:fieldEndPosition]

			column.value = strings.Trim(fieldValue, " ")

			fieldIndex += column.fieldWidth
			record.columns = append(record.columns, column)
		}

		records = append(records, record)
	}

	return records
}

func (f RECFile) ToCSV() {
	csvFilename := f.Filename + ".csv"

	file, err := os.Create(csvFilename)
	check(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	var data = [][]string{}

	var row []string
	for _, column := range f.Header.columns {
		if column.fieldWidth == 0 {
			continue
		}
		row = append(row, column.name[1:])
	}
	data = append(data, row)

	for _, record := range f.Records {
		var row []string
		for _, column := range record.columns {
			row = append(row, column.getStringValue())
		}
		data = append(data, row)
	}

	for _, value := range data {
		err := writer.Write(value)
		check(err)
	}

	fmt.Printf(InfoColor, "[COMPLETE] Conversion REC to CSV")
	fmt.Println("")
}

func (r Column) getStringValue() string {
	//if r.name[0:1] == TEXT {
	//	return r.strValue
	//}
	//return strconv.Itoa(r.intValue)
	return r.value
}
