table 
=====

Table is a Go module providing simple table formatting functions for command line applications.

[![GoDoc](https://godoc.org/github.com/jayloop/table?status.svg)](http://godoc.org/github.com/jayloop/table) [![Go Report](https://goreportcard.com/badge/github.com/jayloop/table)](https://goreportcard.com/report/github.com/jayloop/table)


installation
-------

    go get -u github.com/jayloop/table

usage
-----
```
t := table.New("key", "value")
t.FormatHeader(table.Format(table.HiYellow, table.Bold))
t.Precision(2, 1)
t.Row("a", 1)
t.Row("b", 2.0)
t.Row("c", 0.001)
t.Sort(1)
t.Print(os.Stdout)
```