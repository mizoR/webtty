package webtty

type Option struct {
	InFile string
	Port   int
	Row    int
	Col    int
}

func NewOption() *Option {
	option := &Option{InFile: "ttyrecord", Port: 10101, Row: 24, Col: 120}
	return option
}
