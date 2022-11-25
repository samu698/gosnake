package puppeteer

import (
	. "github.com/samu698/gosnake/snake"
	. "github.com/samu698/gosnake/screen"
)

type Path struct {
	StartingPos Pos
	Path []Direction
}

type Puppeteer struct {
	paths []Path
	snakes []SnakeGame
	frame int
}

const (
	direction_mask Direction = 3
	NO_GROW = 4
)

func NewPuppeteer(screen *Screen, paths []Path) Puppeteer {
	var snakes []SnakeGame
	for _, path := range paths {
		config := SnakeConfig{
			StartDirection: path.Path[0],
			Origin: path.StartingPos,
			CheckBounds: false,
			DrawBorder: false,
			SpawnFruit: false,
			GrowOnFrame: true,
			ClearScreen: false,
		}
		snakes = append(snakes, NewSnakeGame(screen, config))
	}
	return Puppeteer{
		paths: paths,
		snakes: snakes,
		frame: 0,
	}
}

func (this *Puppeteer) Draw() {
	for _, snake := range this.snakes {
		snake.Draw()
	}
}

func (this *Puppeteer) Update() {
	for i, path := range this.paths {
		if this.frame >= len(path.Path) - 1 { continue }
		nextDirection := path.Path[this.frame + 1]
		if nextDirection & NO_GROW != 0 {
			this.snakes[i].Config.GrowOnFrame = false
		}
		this.snakes[i].UpdateWithInput(nextDirection & direction_mask)
		if nextDirection & NO_GROW != 0 {
			this.snakes[i].Config.GrowOnFrame = true
		}
	}
	this.frame++
}
