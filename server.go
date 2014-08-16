package webtty

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"j4k.co/terminal"
)

type Server struct {
	State *terminal.State
	Row   int
	Col   int
}

func NewServer(state *terminal.State, row int, col int) *Server {
	server := &Server{State: state, Row: row, Col: col}

	http.HandleFunc(
		"/",
		staticView("views/index.html"))

	http.HandleFunc(
		"/stylesheets/webtty.css",
		staticView("views/stylesheets/webtty.css"))

	http.HandleFunc(
		"/terminal",
		terminalView(server.State, server.Row, server.Col))

	return server
}

func (server Server) ListenAndServe(port int) {
	log.Printf("== The WebTTY is standing on watch at http://0.0.0.0:%d/", port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", port), Log(http.DefaultServeMux))
	if err != nil {
		panic(err)
	}
}

func staticView(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// FIXME
		t, err := template.ParseFiles(
			filepath.Join(os.Getenv("GOPATH"), "src/github.com/mizoR/webtty", path))
		if err != nil {
			panic(err)
		}
		t.Execute(w, t)
	}
}

func terminalView(state *terminal.State, row int, col int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		writer := NewBufferWriter()
		cx, cy := state.Cursor()

		for c := 0; c < col; c++ {
			for r := 0; r < row; r++ {
				if cx == r && cy == c {
					writer.writeCursor(&buf)
				} else {
					ch, _, _ := state.Cell(r, c)
					writer.write(&buf, ch)
				}
			}
			writer.writeLF(&buf)
		}
		fmt.Fprint(w, buf.String())
	}
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
