// Package table provides simple table formatting functions for command line applications
package table

import (
	"fmt"
	"io"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	defaultHeaderFormat FormatFunc
	mu                  sync.RWMutex
)

// DefaultHeaderFormat sets the default format to use for all new tables.
func DefaultHeaderFormat(f FormatFunc) {
	mu.Lock()
	defaultHeaderFormat = f
	mu.Unlock()
}

// FormatFunc is a user defined function applying formatting to a header or row value.
// The typical usecase is setting colors by adding escape characters.
// A format function should not change the printed length of the value.
type FormatFunc func(string) string

// A Table record stores all table data and formatting options.
// It is not safe for concurrent use.
type Table struct {
	columns       int
	headers       []string
	rows          [][]string
	widths        []int
	maxWidths     []int
	precision     []int
	padding       int
	format        []FormatFunc
	formatHeader  FormatFunc
	formatRow     map[int]FormatFunc
	formatNotZero map[int]FormatFunc
	sortBy        []int
}

// New creates a new table with the given headers.
// The number of headers decides the number of columns of the table.
func New(headers ...string) *Table {
	l := len(headers)
	t := &Table{
		columns:       l,
		headers:       headers,
		widths:        make([]int, l),
		maxWidths:     make([]int, l),
		precision:     make([]int, l),
		format:        make([]FormatFunc, l),
		formatRow:     make(map[int]FormatFunc),
		formatNotZero: make(map[int]FormatFunc),
		rows:          [][]string{},
		padding:       2,
	}
	for i, h := range headers {
		t.widths[i] = len([]rune(h))
	}
	mu.RLock()
	if defaultHeaderFormat != nil {
		t.formatHeader = defaultHeaderFormat
	}
	mu.RUnlock()
	return t
}

// AddStruct adds 2-column rows to a table by iterating over struct fields.
// The table is created by a previous call to New:
//  table.New("key", "value")
func (t *Table) AddStruct(m interface{}) {
	if reflect.TypeOf(m).Kind() != reflect.Struct {
		return
	}
	v := reflect.ValueOf(m)
	for i := 0; i < v.NumField(); i++ {
		t.Row(v.Type().Field(i).Name, v.Field(i).Interface())
	}
}

// FormatHeader sets to format applied to column headers when printing
func (t *Table) FormatHeader(fn FormatFunc) {
	t.formatHeader = fn
}

// Padding sets the number of whitespaces added as padding between columns.
func (t *Table) Padding(p int) {
	t.padding = p
}

// MaxWidth sets the max width in characters for the listed column indexes.
func (t *Table) MaxWidth(chars int, cols ...int) {
	for _, col := range cols {
		if col >= 0 && col < t.columns {
			t.maxWidths[col] = chars
		}
	}
}

// Precision sets the number of digits to include when printing float values.
// It must be set before adding the rows.
func (t *Table) Precision(digits int, cols ...int) {
	for _, col := range cols {
		if col >= 0 && col < t.columns {
			t.precision[col] = digits
		}
	}
}

// FormatRows adds a format function for the listed rows indexes.
// Use row index -1 to denote the last row.
func (t *Table) FormatRows(fn FormatFunc, rows ...int) {
	for _, row := range rows {
		if row == -1 {
			//last row
			row = len(t.rows) - 1
		}
		t.formatRow[row] = fn
	}
}

// FormatCols adds a format function for the listed column indexes.
func (t *Table) FormatCols(fn FormatFunc, cols ...int) {
	for _, col := range cols {
		if col >= 0 && col < t.columns {
			t.format[col] = fn
		}
	}
}

// FormatNotZero adds a format function applied on values != "0"
func (t *Table) FormatNotZero(fn FormatFunc, cols ...int) {
	for _, col := range cols {
		if col >= 0 && col < t.columns {
			t.formatNotZero[col] = fn
		}
	}
}

// Len returns the number of rows
func (t *Table) Len() int {
	return len(t.rows)
}

// Less compares row i against row j
func (t *Table) Less(i, j int) bool {
	var c int
	for _, k := range t.sortBy {
		c = strings.Compare(t.rows[i][k], t.rows[j][k])
		if c != 0 {
			break
		}
	}
	return c < 0
}

// Swap swaps row i and j
func (t *Table) Swap(i, j int) {
	t.rows[i], t.rows[j] = t.rows[j], t.rows[i]
}

// Sort sort the table rows by the listed columns
func (t *Table) Sort(cols ...int) {
	t.sortBy = cols
	sort.Sort(t)
}

// Row adds row data.
func (t *Table) Row(values ...interface{}) {
	// truncate any overflowing values
	if len(values) > t.columns {
		values = values[:t.columns]
	}
	row := make([]string, len(values))
	for i, v := range values {
		p := t.precision[i]
		if p == 0 {
			p = 2
		}
		var v2 string
		switch v := v.(type) {
		case int32:
			v2 = strconv.Itoa(int(v))
		case int64:
			v2 = strconv.FormatInt(v, 10)
		case uint64:
			v2 = strconv.FormatUint(v, 10)
		case float32:
			v2 = strconv.FormatFloat(float64(v), 'f', p, 32)
		case float64:
			v2 = strconv.FormatFloat(v, 'f', p, 64)
		case int:
			v2 = strconv.Itoa(v)
		case uint32:
			v2 = strconv.Itoa(int(v))
		case *[]byte:
			v2 = string(*v)
		case *string:
			v2 = *v
		case nil:
		case bool:
			if v {
				v2 = "yes"
			} else {
				v2 = ""
			}
		case string:
			v2 = v
		default:
			v2 = fmt.Sprintf("%v", v)
		}
		if len([]rune(v2)) > t.widths[i] {
			t.widths[i] = len([]rune(v2))
		}
		row[i] = v2
	}
	t.rows = append(t.rows, row)
}

func appendWhitespace(b []byte, count int) []byte {
	for i := 0; i < count; i++ {
		b = append(b, ' ')
	}
	return b
}

// Print prints the table headers and rows to a io.Writer.
// Any error returned is from the underlying io.Writer.
func (t *Table) Print(out io.Writer) error {
	var buf []byte
	for i, w := range t.widths {
		if t.maxWidths[i] > 0 && w > t.maxWidths[i] {
			t.widths[i] = t.maxWidths[i]
		}
	}
	for i, h := range t.headers {
		if t.maxWidths[i] > 0 && len([]rune(h)) > t.maxWidths[i] {
			h = h[:t.maxWidths[i]-3] + "..."
		}
		l := t.widths[i] + t.padding
		p := l - len([]rune(h))
		if t.formatHeader != nil {
			h = t.formatHeader(h)
		}
		buf = append(buf[:0], []byte(h)...)
		if i != t.columns-1 {
			buf = appendWhitespace(buf, p)
		}
		if _, err := out.Write(buf); err != nil {
			return err
		}
	}
	if _, err := out.Write([]byte("\n")); err != nil {
		return err
	}
	for j, row := range t.rows {
		for i, r := range row {
			if t.maxWidths[i] > 0 && len([]rune(r)) > t.maxWidths[i] {
				r = r[:t.maxWidths[i]-3] + "..."
			}
			l := t.widths[i] + t.padding
			p := l - len([]rune(r))
			switch {
			case t.formatNotZero[i] != nil && r != "0":
				r = t.formatNotZero[i](r)
			case t.formatRow[j] != nil:
				r = t.formatRow[j](r)
			case t.format[i] != nil:
				r = t.format[i](r)
			}
			buf = append(buf[:0], []byte(r)...)
			if i != t.columns-1 {
				buf = appendWhitespace(buf, p)
			}
			if _, err := out.Write(buf); err != nil {
				return err
			}
		}
		if _, err := out.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
