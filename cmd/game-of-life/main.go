package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/johannesjahn/game-of-life/internal/game"
)

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
	tcell.ColorReset,
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

func drawGame(gm game.GameModel) {

	for i, row := range gm.Grid {
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

func writeLine(line string, y, width int) {

	for i := 0; i < width; i++ {
		screen.SetContent(i, y, ' ', nil, defStyle)
	}
	for i, c := range line {
		screen.SetContent(i, y, c, nil, defStyle)
	}
}

var defStyle = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var screen tcell.Screen

type config struct {
	height   int
	width    int
	interval int
	living   int
	seed     int
	factions int
}

func parseInput() config {
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

	return config{
		height:   height,
		width:    width,
		interval: interval,
		living:   living,
		seed:     seed,
		factions: factions,
	}
}

func main() {

	c := parseInput()

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

	if c.width > sw || c.width < 1 {
		c.width = sw
	}
	sh = sh - 2
	if c.height > sh || c.height < 1 {
		c.height = sh
	}
	if c.living < 1 {
		c.living = (c.width * c.height) / 3
	}

	gameModel := game.InitGameModel(c.width, c.height, c.living, c.seed, c.factions)

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

		liveCells := game.CountLiveCells(&gameModel)
		writeLine(fmt.Sprintf("Population: %d Generation: %d", liveCells, generation), c.height+1, c.width)
		s.Show()
		game.GameStep(&gameModel, c.factions)
		generation++
		time.Sleep(time.Duration(c.interval) * time.Millisecond)
	}
}
