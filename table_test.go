package table_test

import (
	"os"
	"runtime"
	"testing"

	"github.com/jayloop/table"
)

func Example() {
	t := table.New("key", "value")
	t.Precision(2, 1)
	t.Row("a", 1)
	t.Row("b", 2.0)
	t.Row("c", 0.001)
	t.Sort(1)
	t.Print(os.Stdout)
	// Output:
	// key  value
	// c    0.00
	// a    1
	// b    2.00
}

func TestTable(t *testing.T) {
	tbl := table.New("key", "value", "other")
	tbl.FormatHeader(table.Format(table.HiYellow, table.Background(table.Black), table.Bold, table.Underline))
	tbl.FormatCols(table.Format(table.White, table.Background(table.Blue)), 0)
	tbl.FormatCols(table.Format(table.Background(table.Green)), 1)
	tbl.Precision(4, 1)
	tbl.Precision(8, 2)
	tbl.Row(int(1), float64(1), float64(2))
	s := "hello"
	tbl.Row(int(1), &s, 0.5)
	tbl.Row(uint64(0), []int{1, 2, 3, 4})
	tbl.Sort(0, 2)
	tbl.Print(os.Stdout)
}

func TestTableFromStruct(t *testing.T) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	tbl := table.New("key", "value")
	tbl.FormatHeader(table.Format(table.HiYellow, table.Bold))
	tbl.FormatCols(table.Format(table.HiCyan), 0)
	tbl.Precision(4, 1)
	tbl.MaxWidth(60, 1)
	tbl.AddStruct(m)
	tbl.Sort(0)
	tbl.Print(os.Stdout)
}
