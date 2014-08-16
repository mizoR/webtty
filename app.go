package webtty

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"j4k.co/terminal"
)

type App struct {
	State terminal.State
}

func NewApp() *App {
	app := &App{State: terminal.State{}}
	return app
}

func (app App) Run(option *Option) error {
	inFile := "ttyrecord"
	port := 10101
	row := option.Row
	col := option.Col
	state := &app.State

	player := NewPlayer(state, row, col)
	player.Play(inFile)

	http.HandleFunc(
		"/",
		staticView("views/index.html"))

	http.HandleFunc(
		"/stylesheets/webtty.css",
		staticView("views/stylesheets/webtty.css"))

	http.HandleFunc(
		"/terminal",
		terminalView(state, row, col))

	log.Printf("== The WebTTY is standing on watch at http://0.0.0.0:%d/", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), Log(http.DefaultServeMux))
	if err != nil {
		panic(err)
	}
	return nil
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func staticView(filepath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(filepath)
		if err != nil {
			panic(err)
		}
		t.Execute(w, t)
	}
}

func terminalView(state *terminal.State, row int, col int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		cx, cy := state.Cursor()

		for c := 0; c < col; c++ {
			for r := 0; r < row; r++ {
				if cx == r && cy == c {
					writeCursor(&buf)
				} else {
					ch, _, _ := state.Cell(r, c)
					writeRuneAsSecureHTML(&buf, ch)
				}
			}
			writeLF(&buf)
		}
		fmt.Fprint(w, buf.String())
	}
}

func writeLF(buf *bytes.Buffer) {
	buf.WriteRune(10) // LF
}

func writeCursor(buf *bytes.Buffer) {
	buf.WriteString("<div class='cursor'></div>")
}

func writeRuneAsSecureHTML(buf *bytes.Buffer, r rune) {
	switch r {
	case 34: // `"`
		buf.WriteString("&quot;")
	case 38: // `&`
		buf.WriteString("&amp;")
	case 39: // `'`
		buf.WriteString("&#039;")
	case 60: // `<`
		buf.WriteString("&lt;")
	case 62: // `>`
		buf.WriteString("&gt;")
	default:
		buf.WriteRune(r)
	}
}
