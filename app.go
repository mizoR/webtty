package webtty

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/sugyan/ttyread"
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
	port := 10101
	row := option.Row
	col := option.Col

	vt, err := terminal.Create(&(app.State), ioutil.NopCloser(bytes.NewBuffer([]byte{})))
	if err != nil {
		panic(err)
	}
	defer vt.Close()

	vt.Resize(row, col)

	in, err := os.Open("ttyrecord")
	if err != nil {
		panic(err)
	}
	reader := ttyread.NewTtyReader(in)

	go play(vt, reader)

	http.HandleFunc(
		"/",
		staticView("views/index.html"))

	http.HandleFunc(
		"/stylesheets/webtty.css",
		staticView("views/stylesheets/webtty.css"))

	http.HandleFunc(
		"/terminal",
		terminalView(&(app.State), row, col))

	log.Printf("== The WebTTY is standing on watch at http://0.0.0.0:%d/", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), Log(http.DefaultServeMux))
	if err != nil {
		panic(err)
	}
	return nil
}

func play(vt *terminal.VT, reader *ttyread.TtyReader) {
	for {
		data, err := reader.ReadData()
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				panic(err)
			}
		}
		_, err = vt.Write(*data.Buffer)
		time.Sleep(100000000)
	}
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
