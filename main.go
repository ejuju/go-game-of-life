package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

const Width, Height = 40, 40

type Cell bool
type Grid [Height][Width]Cell

func (g Grid) At(x, y int) Cell { return g[(Height+y)%Height][(Width+x)%Width] }

func (g Grid) ForEach(callback func(x, y int, c Cell)) {
	for y, row := range g {
		for x := range row {
			callback(x, y, g.At(x, y))
		}
	}
}

func NewRandomGrid(seed int64) Grid {
	random := rand.New(rand.NewSource(seed))
	g := Grid{}
	g.ForEach(func(x, y int, c Cell) { g[y][x] = random.Intn(8) == 0 })
	return g
}

func (g Grid) Next() Grid {
	next := Grid{}
	g.ForEach(func(x, y int, c Cell) {
		count := g.CountNeighbours(x, y)
		if c {
			switch {
			case count == 2 || count == 3:
				next[y][x] = true
			}
		} else if !c && count == 3 {
			next[y][x] = true
		}
	})
	return next
}

func (g Grid) CountNeighbours(x, y int) int {
	count := 0
	for i := y - 1; i <= y+1; i++ {
		for j := x - 1; j <= x+1; j++ {
			if i == y && j == x {
				continue
			}
			if g.At(j, i) {
				count++
			}
		}
	}
	return count
}

func (g Grid) Render() {
	g.ForEach(func(x, y int, c Cell) {
		line, column := y+1, x*2+1
		print(fmt.Sprintf("\x1b[%d;%dH", line, column)) //  move to position
		if c {
			print("\x1b[42m" + fmt.Sprintf("%d ", g.CountNeighbours(x, y)) + "\x1b[0m")
		} else {
			print("\x1b[40m" + "  " + "\x1b[0m")
		}
	})
}

func main() {
	print("\x1b[?25l")       // Hide cursor
	defer print("\x1b[?25h") // Make cursor visible again on exit

	const fps = 3
	go func() {
		grid := NewRandomGrid(time.Now().Unix())
		for range time.Tick(time.Second / time.Duration(fps)) {
			grid.Render()
			grid = grid.Next()
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)
	<-exit
}

func print(s string) {
	_, err := io.WriteString(os.Stdout, s)
	if err != nil {
		panic(err)
	}
}
