package webtty

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/sugyan/ttyread"
	"j4k.co/terminal"
)

type Player struct {
	State *terminal.State
	Row   int
	Col   int
}

func NewPlayer(state *terminal.State, row int, col int) *Player {
	player := &Player{State: state, Row: row, Col: col}
	return player
}

func (player Player) Play(inFile string) {
	go play(player.State, player.Row, player.Col, inFile)
}

func play(state *terminal.State, row int, col int, inFile string) {
	in, err := os.Open(inFile)
	if err != nil {
		panic(err)
	}
	reader := ttyread.NewTtyReader(in)

	vt, err := terminal.Create(state, ioutil.NopCloser(bytes.NewBuffer([]byte{})))
	if err != nil {
		panic(err)
	}
	defer vt.Close()

	vt.Resize(row, col)

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
