package main

import (
	"flag"

	"github.com/mizoR/webtty"
)

func main() {
	inFile := flag.String("in", "ttyrecord", "in file")
	port := flag.Int("port", 10101, "port")
	row := flag.Int("row", 120, "rows")
	col := flag.Int("col", 24, "columns")
	flag.Parse()

	option := webtty.NewOption()
	option.InFile = *inFile
	option.Port = *port
	option.Row = *row
	option.Col = *col

	app := webtty.NewApp()
	app.Run(option)
}
