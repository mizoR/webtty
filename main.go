package main

import (
	"flag"

	"./webtty"
)

func main() {
	row := flag.Int("row", 100, "rows")
	col := flag.Int("col", 24, "columns")
	flag.Parse()

	option := webtty.NewOption()
	option.Row = *row
	option.Col = *col

	app := webtty.NewApp()
	app.Run(option)
}
