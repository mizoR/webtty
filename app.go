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

	// Play
	in, err := os.Open("ttyrecord")
	reader := ttyread.NewTtyReader(in)

	go func() error {
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

		return nil
	}()

	// http
	http.HandleFunc(
		"/",
		staticViewHandler("views/index.html"))

	http.HandleFunc(
		"/stylesheets/webtty.css",
		staticViewHandler("views/stylesheets/webtty.css"))

	http.HandleFunc("/terminal", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		cx, cy := app.State.Cursor()
		for c := 0; c < col; c++ {
			for r := 0; r < row; r++ {
				if cx == r && cy == c {
					buf.WriteString("<div class='cursor'></div>")
				} else {
					ch, _, _ := app.State.Cell(r, c)

					// HTML escape
					switch ch {
					case 34:
						buf.WriteString("&quot;")
					case 38:
						buf.WriteString("&amp;")
					case 39:
						buf.WriteString("&#039;")
					case 60:
						buf.WriteString("&lt;")
					case 62:
						buf.WriteString("&gt;")
					default:
						buf.WriteRune(ch)
					}
				}
			}
			buf.WriteString("\n")
		}
		fmt.Fprint(w, buf.String())
	})

	log.Printf("== The WebTTY is standing on watch at http://0.0.0.0:%d/", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), Log(http.DefaultServeMux))
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

func staticViewHandler(filepath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles(filepath)
		if err != nil {
			panic(err)
		}
		t.Execute(w, t)
	}
}
