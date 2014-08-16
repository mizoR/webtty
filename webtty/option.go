package webtty

type Option struct {
	Row int
	Col int
}

func NewOption() *Option {
	option := &Option{Row: 24, Col: 120}
	return option
}
