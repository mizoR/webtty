package webtty

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
	row := option.Row
	col := option.Col

	vt, err := terminal.Create(&(app.State), ioutil.NopCloser(bytes.NewBuffer([]byte{})))
	if err != nil {
		return err
	}
	defer vt.Close()

	vt.Resize(row, col)
	vt.Write([]byte("ls\n"))

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
					return err
				}
			}
			_, err = vt.Write(*data.Buffer)
			time.Sleep(100000000)
		}

		return nil
	}()

	// http
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html>
				<head>
					<title>Terminal</title>
				</head>
				<body>
				<h1>Terminal</h1>
				<pre id="terminal">
				</pre>
				<script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
				<script type="text/javascript">
					var callback = function() {
						$('#terminal').load('/terminal', function(text, status, xhr) {
							if (status === 'success') {
								setTimeout(callback, 100);
							}
						})
					}
					setTimeout(callback, 50);
				</script>
				</body>
			</html>
		`)
	})

	http.HandleFunc("/terminal", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		for c := 0; c < col; c++ {
			for r := 0; r < row; r++ {
				ch, _, _ := app.State.Cell(r, c)
				buf.WriteRune(ch)
			}
			buf.WriteString("\n")
		}
		fmt.Fprintf(w, "%s", buf.String())
	})

	err = http.ListenAndServe(":10101", nil)
	if err != nil {
		return err
	}

	return nil
}
