package game

import (
	"log"
	"math/rand"
)

type GameModel struct {
	Grid [][]rune
}

func InitGameModel(width, height, living, seed, factions int) GameModel {
	log.Println("Initializing game model with width:", width, "height:", height, "living:", living, "seed:", seed)
	gm := GameModel{
		Grid: make([][]rune, height),
	}
	for i := range gm.Grid {
		gm.Grid[i] = make([]rune, width)
		for j := range gm.Grid[i] {
			gm.Grid[i][j] = '.'
		}
	}

	if living > width*height {
		living = width * height
	}

	r := rand.New(rand.NewSource(int64(seed)))
	for living > 0 {

		x := r.Intn(height)
		y := r.Intn(width)
		if gm.Grid[x][y] == '.' {
			if factions < 2 {
				gm.Grid[x][y] = 'O'
			} else {
				gm.Grid[x][y] = rune(r.Intn(factions) + 'A')
			}
			living--
		}
	}

	log.Println("Game model initialized")
	return gm
}

func GameStep(gm *GameModel, factions int) {
	nextGrid := make([][]rune, len(gm.Grid))
	for i := range gm.Grid {
		nextGrid[i] = make([]rune, len(gm.Grid[i]))
		for j := range gm.Grid[i] {
			kind := gm.Grid[i][j]
			if kind == '.' {
				if factions < 2 {
					kind = 'O'
				} else {
					kind = rune('A' + rand.Intn(factions))
				}
			}
			liveNeighbors := CountLiveNeighbors(gm, i, j, kind)
			if gm.Grid[i][j] == kind {
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
	gm.Grid = nextGrid
}

func CountLiveNeighbors(gm *GameModel, x, y int, kind rune) int {
	directions := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	count := 0
	for _, d := range directions {
		nx, ny := x+d.dx, y+d.dy
		if nx >= 0 && nx < len(gm.Grid) && ny >= 0 && ny < len(gm.Grid[0]) && gm.Grid[nx][ny] == kind {
			count++
		}
	}
	return count
}

func CountLiveCells(gm *GameModel) int {
	count := 0
	for i := range gm.Grid {
		for j := range gm.Grid[i] {
			if gm.Grid[i][j] != '.' {
				count++
			}
		}
	}
	return count
}
