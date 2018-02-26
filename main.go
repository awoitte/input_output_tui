package input_output_tui

import (
	"fmt"
	"os"
	"time"

	"github.com/gdamore/tcell"
)

type State struct {
	input  string
	output []string
}

func redraw_screen(s tcell.Screen, state *State) {
	w, h := s.Size()

	if w == 0 || h == 0 {
		return
	}

	st := tcell.StyleDefault

	for row := 0; row < h; row++ {
		var line_text string
		from_bottom := (h - 1) - row
		if from_bottom == 0 {
			line_text = state.input
		} else if from_bottom <= len(state.output) {
			line_text = state.output[len(state.output)-from_bottom]
		} else {
			line_text = ""
		}

		for col := 0; col < w; col++ {

			var letter byte
			if len(line_text) > col {
				letter = line_text[col]
			} else {
				letter = ' '
			}
			s.SetCell(col, row, st, rune(letter))

		}
	}
	s.Show()
}

func respond_to_input(s tcell.Screen, state *State, quit chan bool, user_input chan string) {
	for {
		ev := s.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			key := ev.Key()
			switch key {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				close(quit)
				return
			case tcell.KeyCtrlL:
				s.Sync()
			case tcell.KeyRune:
				state.input = state.input + string(ev.Rune())
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(state.input) > 0 {
					state.input = state.input[:len(state.input)-1]
				}
			case tcell.KeyEnter:
				user_input <- state.input
				state.input = ""
			}
			redraw_screen(s, state)
		case *tcell.EventResize:
			s.Sync()
		}
	}
}

func Start(user_input chan string, output_messages chan string, quit chan bool) {
	state := State{"", []string{}}

	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	defer s.Fini()

	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorWhite).
		Background(tcell.ColorBlack))
	s.Clear()

	go respond_to_input(s, &state, quit, user_input)

	for {
		select {
		case <-quit:
			return
		case message := <-output_messages:
			state.output = append(state.output, message)
		case <-time.After(time.Second):
		}
		redraw_screen(s, &state)
	}
}
