package table

import (
	"os"
	"runtime"
	"testing"
)

func Example() {
	t := New("key", "value")
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
	tbl := New("key", "value", "other")
	tbl.FormatHeader(Format(HiYellow, Background(Black), Bold, Underline))
	tbl.FormatCols(Format(White, Background(Blue)), 0)
	tbl.FormatCols(Format(Background(Green)), 1)
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
	tbl := New("key", "value")
	tbl.FormatHeader(Format(HiYellow, Bold))
	tbl.FormatCols(Format(HiCyan), 0)
	tbl.Precision(4, 1)
	tbl.MaxWidth(60, 1)
	tbl.AddStruct(m)
	tbl.Sort(0)
	tbl.Print(os.Stdout)
}
