package main

import (
	"fmt"
	"os"

	"github.com/gdamore/tcell"
)

type Position struct {
	row int
	col int
}

const (
	NORMAL  = iota
	INSERT  = iota
	REPLACE = iota
	VISUAL  = iota
)

func main() {
	scr, err := tcell.NewScreen()
	checkError(err)
	defer quit(scr)
	err = scr.Init()
	checkError(err)
	normalStyle := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	insertStyle := tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack).Blink(true)
	mode := NORMAL
	style := normalStyle
	lines := []string{""}
	col := 0
	row := 0
	for {
		if mode == NORMAL {
			style = normalStyle
		} else if mode == INSERT {
			style = insertStyle
		}
		scr.Clear()
		cols, rows := scr.Size()
		for r, l := range lines {
			drawText(scr, 0, r, style, l)
		}
		scr.ShowCursor(col, row)
		scr.Show()
		ev := scr.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			scr.Sync()
		case *tcell.EventKey:
			if mode == NORMAL {
				if ev.Key() == tcell.KeyCtrlC {
					quit(scr)
				}
				if ev.Rune() == 'h' {
					col -= 1
				}
				if ev.Rune() == 'j' {
					row += 1
					if len(lines) <= row {
						lines = append(lines, "")
					}
					col = min(len(lines[row]), col)
				}
				if ev.Rune() == 'k' {
					row -= 1
					col = min(len(lines[row]), col)
				}
				if ev.Rune() == 'l' {
					col += 1
				}
				if ev.Rune() == 'i' {
					mode = INSERT
					col = len(lines[row]) - 1
				}
			} else if mode == INSERT {
				if ev.Key() == tcell.KeyEscape {
					mode = NORMAL
				} else if ev.Key() == tcell.KeyEnter {
					lines = append(lines, "")
					row++
					col = 0
				} else if ev.Key() == tcell.KeyBackspace2 {
					if len(lines[row]) > 0 {
						lines[row] = lines[row][0 : len(lines[row])-1]
					}
					col -= 1
					if col < 0 && row > 0 {
						row--
						col = len(lines[row])
					}
				} else {
					lines[row] += string(ev.Rune())
					col += 1
				}
			}
			if col < 0 {
				col = 0
			}
			if col >= cols {
				col = cols - 1
			}
			if row < 0 {
				row = 0
			}
			if row >= rows {
				row = rows - 1
			}
		}
	}
}

func drawText(scr tcell.Screen, col int, row int, style tcell.Style, text string) {
	for _, r := range []rune(text) {
		scr.SetContent(col, row, r, nil, style)
		col++
	}
}

func drawTextWrapping(scr tcell.Screen, startCol int, startRow int, endCol int, endRow int, style tcell.Style, text string) {
	col := startCol
	row := startRow
	for _, r := range []rune(text) {
		if r == '\n' {
			row++
			col = startCol
		} else {
			scr.SetContent(col, row, r, nil, style)
			col++
		}
		if col >= endCol {
			row++
			col = startCol
		}
		if row > endRow {
			break
		}
	}
}

func quit(scr tcell.Screen) {
	scr.Fini()
	os.Exit(0)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
