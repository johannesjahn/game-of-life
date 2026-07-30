package game

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
)

// directions is the fixed 8-neighbor offset table, built once instead of on
// every CountLiveNeighbors call.
var directions = [8]struct{ dx, dy int }{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

type GameModel struct {
	Width  int
	Height int
	Grid   []rune

	next []rune
	rng  *rand.Rand
}

func InitGameModel(width, height, living, seed, factions int) GameModel {
	log.Println("Initializing game model with width:", width, "height:", height, "living:", living, "seed:", seed)
	gm := GameModel{
		Width:  width,
		Height: height,
		Grid:   make([]rune, width*height),
		next:   make([]rune, width*height),
		rng:    rand.New(rand.NewSource(int64(seed))),
	}
	for i := range gm.Grid {
		gm.Grid[i] = '.'
	}

	if living > width*height {
		living = width * height
	}

	r := gm.rng
	for living > 0 {
		x := r.Intn(height)
		y := r.Intn(width)
		idx := x*width + y
		if gm.Grid[idx] == '.' {
			if factions < 2 {
				gm.Grid[idx] = 'O'
			} else {
				gm.Grid[idx] = rune(r.Intn(factions) + 'A')
			}
			living--
		}
	}

	log.Println("Game model initialized")
	return gm
}

func GameStep(gm *GameModel, factions int) {
	height := gm.Height

	numWorkers := runtime.GOMAXPROCS(0)
	if numWorkers > height {
		numWorkers = height
	}
	if numWorkers < 1 {
		numWorkers = 1
	}
	rowsPerWorker := (height + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	for w := range numWorkers {
		startRow := w * rowsPerWorker
		endRow := startRow + rowsPerWorker
		if endRow > height {
			endRow = height
		}
		if startRow >= endRow {
			break
		}

		// Each worker gets its own RNG, seeded deterministically from the
		// model's seeded source, so no shared/global rand lock is hit on
		// the hot path and results stay reproducible for a given seed.
		workerRand := rand.New(rand.NewSource(gm.rng.Int63()))

		wg.Add(1)
		go func(startRow, endRow int, r *rand.Rand) {
			defer wg.Done()
			stepRows(gm, r, factions, startRow, endRow)
		}(startRow, endRow, workerRand)
	}
	wg.Wait()

	gm.Grid, gm.next = gm.next, gm.Grid
}

func stepRows(gm *GameModel, r *rand.Rand, factions, startRow, endRow int) {
	width := gm.Width
	for i := startRow; i < endRow; i++ {
		base := i * width
		for j := range width {
			idx := base + j
			kind := gm.Grid[idx]
			if kind == '.' {
				if factions < 2 {
					kind = 'O'
				} else {
					kind = rune('A' + r.Intn(factions))
				}
			}
			liveNeighbors := countLiveNeighbors(gm, i, j, kind)
			if gm.Grid[idx] == kind {
				if liveNeighbors == 2 || liveNeighbors == 3 {
					gm.next[idx] = kind
				} else {
					gm.next[idx] = '.'
				}
			} else {
				if liveNeighbors == 3 {
					gm.next[idx] = kind
				} else {
					gm.next[idx] = '.'
				}
			}
		}
	}
}

func countLiveNeighbors(gm *GameModel, x, y int, kind rune) int {
	width, height := gm.Width, gm.Height
	count := 0
	for _, d := range directions {
		nx, ny := x+d.dx, y+d.dy
		if nx >= 0 && nx < height && ny >= 0 && ny < width && gm.Grid[nx*width+ny] == kind {
			count++
		}
	}
	return count
}

func CountLiveCells(gm *GameModel) int {
	count := 0
	for _, cell := range gm.Grid {
		if cell != '.' {
			count++
		}
	}
	return count
}
