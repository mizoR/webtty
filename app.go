package webtty

import "j4k.co/terminal"

type App struct {
	State terminal.State
}

func NewApp() *App {
	app := &App{State: terminal.State{}}
	return app
}

func (app App) Run(option *Option) {
	inFile := option.InFile
	port := option.Port
	row := option.Row
	col := option.Col
	state := &app.State

	player := NewPlayer(state, row, col)
	player.Play(inFile)

	server := NewServer(state, row, col)
	server.ListenAndServe(port)
}
