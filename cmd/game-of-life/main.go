package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type gameModel struct {
	grid [][]rune
}

func drawGame(gm gameModel) {

	for i, row := range gm.grid {
		for j, cell := range row {
			screen.SetContent(j, i, cell, nil, defStyle)
		}
	}
}

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var screen tcell.Screen

func initGameModel(width, height int) gameModel {
	gm := gameModel{
		grid: make([][]rune, height),
	}
	for i := range gm.grid {
		gm.grid[i] = make([]rune, width)
		for j := range gm.grid[i] {
			gm.grid[i][j] = '.'
		}
	}

	gm.grid[5][5] = 'O'
	gm.grid[5][6] = 'O'
	gm.grid[5][7] = 'O'
	gm.grid[6][5] = 'O'
	gm.grid[6][6] = 'O'
	gm.grid[6][7] = 'O'
	gm.grid[7][5] = 'O'
	gm.grid[7][6] = 'O'
	gm.grid[7][7] = 'O'
	gm.grid[8][5] = 'O'

	return gm
}

func gameStep(gm *gameModel) {
	nextGrid := make([][]rune, len(gm.grid))
	for i := range gm.grid {
		nextGrid[i] = make([]rune, len(gm.grid[i]))
		for j := range gm.grid[i] {
			liveNeighbors := countLiveNeighbors(gm, i, j)
			if gm.grid[i][j] == 'O' {
				if liveNeighbors < 2 || liveNeighbors > 3 {
					nextGrid[i][j] = '.'
				} else {
					nextGrid[i][j] = 'O'
				}
			} else {
				if liveNeighbors == 3 {
					nextGrid[i][j] = 'O'
				} else {
					nextGrid[i][j] = '.'
				}
			}
		}
	}
	gm.grid = nextGrid
}

func countLiveNeighbors(gm *gameModel, x, y int) int {
	directions := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	count := 0
	for _, d := range directions {
		nx, ny := x+d.dx, y+d.dy
		if nx >= 0 && nx < len(gm.grid) && ny >= 0 && ny < len(gm.grid[0]) && gm.grid[nx][ny] == 'O' {
			count++
		}
	}
	return count
}

func main() {

	var (
		height int
		width  int
	)

	flag.IntVar(&height, "height", 10, "height of the grid (default 10)")
	flag.IntVar(&height, "h", 10, "height of the grid (default 10) (shorthand)")
	flag.IntVar(&width, "width", 10, "width of the grid (default 10)")
	flag.IntVar(&width, "w", 10, "width of the grid (default 10) (shorthand)")

	// Parse the flags
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Println("Positional arguments found")
		os.Exit(2)
	}

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	screen = s
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	gameModel := initGameModel(width, height)

	go func() {
		for {

			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape {
					s.Fini()
					return
				}
			}
		}
	}()

	for {
		// Update screen
		drawGame(gameModel)
		s.Show()
		gameStep(&gameModel)

		// Poll for events

		time.Sleep(1000 * time.Millisecond)
	}
}
