package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

type gameModel struct {
	grid [][]rune
}

var colors = []tcell.Color{
	tcell.ColorRed,
	tcell.ColorGreen,
	tcell.ColorBlue,
	tcell.ColorYellow,
	tcell.ColorDarkRed,
	tcell.ColorDarkGreen,
	tcell.ColorDarkBlue,
	tcell.ColorDarkGoldenrod,
	tcell.ColorOrangeRed,
	tcell.ColorDarkSlateGray,
	tcell.ColorDarkOliveGreen,
	tcell.ColorDarkOrchid,
	tcell.ColorDarkSalmon,
	tcell.ColorDarkSeaGreen,
	tcell.ColorDarkTurquoise,
	tcell.ColorDarkViolet,
	tcell.ColorDeepPink,
	tcell.ColorDeepSkyBlue,
	tcell.ColorDimGray,
	tcell.ColorDodgerBlue,
	tcell.ColorFireBrick,
	tcell.ColorFloralWhite,
	tcell.ColorForestGreen,
	tcell.ColorFuchsia,
	tcell.ColorGainsboro,
}

func drawGame(gm gameModel) {

	for i, row := range gm.grid {
		for j, cell := range row {
			if cell != '.' {
				color := colors[(int(cell)-'A')%len(colors)]
				style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(color)
				screen.SetContent(j, i, cell, nil, style)
			} else {
				screen.SetContent(j, i, cell, nil, defStyle)
			}
		}
	}
}

func writeLine(line string, y int) {
	_, width := screen.Size()
	for i := 0; i < width; i++ {
		screen.SetContent(i, y, ' ', nil, defStyle)
	}
	for i, c := range line {
		screen.SetContent(i, y, ' ', nil, defStyle)
		screen.SetContent(i, y, c, nil, defStyle)
	}
}

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var screen tcell.Screen

func initGameModel(width, height, living, seed, factions int) gameModel {
	log.Println("Initializing game model with width:", width, "height:", height, "living:", living, "seed:", seed)
	gm := gameModel{
		grid: make([][]rune, height),
	}
	for i := range gm.grid {
		gm.grid[i] = make([]rune, width)
		for j := range gm.grid[i] {
			gm.grid[i][j] = '.'
		}
	}

	if living > width*height {
		living = width * height
	}

	r := rand.New(rand.NewSource(int64(seed)))
	for living > 0 {

		x := r.Intn(height)
		y := r.Intn(width)
		if gm.grid[x][y] == '.' {
			if factions < 2 {
				gm.grid[x][y] = 'O'
			} else {
				gm.grid[x][y] = rune(r.Intn(factions) + 'A')
			}
			living--
		}
	}

	log.Println("Game model initialized")
	return gm
}

func gameStep(gm *gameModel, factions int) {
	nextGrid := make([][]rune, len(gm.grid))
	for i := range gm.grid {
		nextGrid[i] = make([]rune, len(gm.grid[i]))
		for j := range gm.grid[i] {
			kind := gm.grid[i][j]
			if kind == '.' {
				if factions < 2 {
					kind = 'O'
				} else {
					kind = rune('A' + rand.Intn(factions))
				}
			}
			liveNeighbors := countLiveNeighbors(gm, i, j, kind)
			if gm.grid[i][j] == kind {
				if liveNeighbors == 2 || liveNeighbors == 3 {
					nextGrid[i][j] = kind
				} else {
					nextGrid[i][j] = '.'
				}
			} else {
				if liveNeighbors == 3 {
					nextGrid[i][j] = kind
				} else {
					nextGrid[i][j] = '.'
				}
			}
		}
	}
	gm.grid = nextGrid
}

func countLiveNeighbors(gm *gameModel, x, y int, kind rune) int {
	directions := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	count := 0
	for _, d := range directions {
		nx, ny := x+d.dx, y+d.dy
		if nx >= 0 && nx < len(gm.grid) && ny >= 0 && ny < len(gm.grid[0]) && gm.grid[nx][ny] == kind {
			count++
		}
	}
	return count
}

func countLiveCells(gm *gameModel) int {
	count := 0
	for i := range gm.grid {
		for j := range gm.grid[i] {
			if gm.grid[i][j] != '.' {
				count++
			}
		}
	}
	return count
}

func main() {

	var (
		height   int
		width    int
		interval int
		living   int
		seed     int
		factions int
	)

	flag.IntVar(&height, "height", -1, "height of the grid (default max possible)")
	flag.IntVar(&height, "h", -1, "height of the grid (default max possible) (shorthand)")
	flag.IntVar(&width, "width", -1, "width of the grid (default max possible)")
	flag.IntVar(&width, "w", -1, "width of the grid (default max possible) (shorthand)")
	flag.IntVar(&interval, "interval", 100, "interval between steps in milliseconds (default 100)")
	flag.IntVar(&interval, "i", 100, "interval between steps in milliseconds (default 100) (shorthand)")
	flag.IntVar(&living, "living", -1, "number of living cells (default (width * height) / 3)")
	flag.IntVar(&living, "l", -1, "number of living cells (default (width * height) / 3) (shorthand)")
	flag.IntVar(&seed, "seed", 0, "seed for random number generator (default 0)")
	flag.IntVar(&seed, "s", 0, "seed for random number generator (default 0) (shorthand)")
	flag.IntVar(&factions, "factions", 0, "number of factions (default 0)")
	flag.IntVar(&factions, "f", 0, "number of factions (default 0) (shorthand)")

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

	sw, sh := s.Size()

	if width > sw || width < 1 {
		width = sw
	}
	sh = sh - 2
	if height > sh || height < 1 {
		height = sh
	}
	if living < 1 {
		living = (width * height) / 3
	}

	gameModel := initGameModel(width, height, living, seed, factions)

	exitChan := make(chan struct{})

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Key() == tcell.KeyCtrlC {
					close(exitChan)
					return
				}
			}
		}
	}()

	generation := 0
	for {
		select {
		case <-exitChan:
			return
		default:
		}
		// Update screen
		drawGame(gameModel)

		liveCells := countLiveCells(&gameModel)
		writeLine(fmt.Sprintf("Population: %d Generation: %d", liveCells, generation), height+1)
		s.Show()
		gameStep(&gameModel, factions)
		generation++
		time.Sleep(time.Duration(interval) * time.Millisecond)
	}
}
