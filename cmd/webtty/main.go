package main

import (
	"flag"

	"github.com/mizoR/webtty/webtty"
)

func main() {
	row := flag.Int("row", 120, "rows")
	col := flag.Int("col", 24, "columns")
	flag.Parse()

	option := webtty.NewOption()
	option.Row = *row
	option.Col = *col

	app := webtty.NewApp()
	app.Run(option)
}
