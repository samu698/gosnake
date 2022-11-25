package snake

import (
	"math/rand"
	"strings"
	"time"

	. "github.com/samu698/gosnake/screen"
)

type Direction int
const (
	LEFT Direction = iota
	RIGHT
	UP
	DOWN
)

type SnakeGame struct {
	myRand *rand.Rand
	screen *Screen
	tailDirection Direction
	snake []snakeSegment
	fruit Pos
	Config SnakeConfig
}

type SnakeConfig struct {
	StartDirection Direction
	Size, Origin Pos
	CheckBounds bool
	DrawBorder bool
	SpawnFruit bool
	GrowOnFrame bool
	ClearScreen bool
}

// map[myDir][prevDir]
var snakeCharset = map[Direction]map[Direction]string {
	LEFT: {
		LEFT: "══",
		UP: "╗ ",
		DOWN: "╝ ",
	},
	RIGHT: {
		RIGHT: "══",
		UP: "╔═",
		DOWN: "╚═",
	},
	UP: {
		UP: "║ ",
		LEFT: "╚═",
		RIGHT: "╝ ",
	},
	DOWN: {
		DOWN: "║ ",
		LEFT: "╔═",
		RIGHT: "╗ ",
	},
}

type snakeSegment struct {
	pos Pos
	direction Direction
}

func NewSnakeGame(screen *Screen, config SnakeConfig) SnakeGame {
	source := rand.NewSource(time.Now().UnixNano())
	myRand := rand.New(source)

	snake := make([]snakeSegment, 1)
	snake[0] = snakeSegment{config.Origin, config.StartDirection}

	var fruit Pos
	if config.SpawnFruit {
		fruit = NewPos(myRand.Intn(config.Size.X), myRand.Intn(config.Size.Y))
	} else {
		fruit = NewPos(-1, -1)
	}

	screen.StartInputReading()

	return SnakeGame{
		myRand: myRand,
		screen: screen,
		tailDirection: config.StartDirection,
		snake: snake,
		fruit: fruit,
		Config: config,
	}
}

func (this *SnakeGame) randomFruit() (fruit Pos) {
	outer:
	for {
		fruit = NewPos(this.myRand.Intn(this.Config.Size.X), this.myRand.Intn(this.Config.Size.Y))
		for _, s := range this.snake {
			if fruit.Equal(s.pos) { continue outer }
		}
		return fruit
	}
}

func (this *SnakeGame) UpdateWithInput(direction Direction) (dead bool) {
	this.snake[0].direction = direction

	head := this.snake[0]

	switch head.direction {
	case LEFT: head.pos.Addv(-1, 0)
	case RIGHT: head.pos.Addv(+1, 0)
	case UP: head.pos.Addv(0, -1)
	case DOWN: head.pos.Addv(0, +1)
	}

	if this.Config.GrowOnFrame || head.pos.Equal(this.fruit) {
		if this.Config.SpawnFruit {
			this.fruit = this.randomFruit()
		}

		// If the snake has grown
		// we can simply insert the new head into the slice
		this.snake = append(this.snake, snakeSegment{NewPos(0, 0), head.direction})
		copy(this.snake[1:], this.snake[0:])
	} else {
		// Save the old tail direction,
		// so that the tail can be drawn correctly
		this.tailDirection = this.snake[len(this.snake)-1].direction

		// Move the snake, by traslating all the positions
		// and discarding the last one
		if len(this.snake) > 1 {
			copy(this.snake[1:], this.snake[0:])
		}

		// Check if the snake collided with itself
		// This check is not done when the snake eats the fruit
		// because they cannot spawn into the snake
		if this.Config.CheckBounds {
			for _, v := range this.snake {
				if v.pos.Equal(head.pos) {
					dead = true
					return
				}
			}
		}
	}

	this.snake[0] = head

	dead = !head.pos.IsInside(NewPos(0, 0), NewPos(this.Config.Size.X - 1, this.Config.Size.Y - 1))
	return

}
func (this *SnakeGame) Update() (dead bool) {
	oldDirection := this.snake[0].direction
	newDirection := oldDirection

	this.screen.ReadInput(func(k Keycode) {
		switch k {
		case Up: if oldDirection != DOWN { newDirection = UP }
		case Down: if oldDirection != UP { newDirection = DOWN }
		case Left: if oldDirection != RIGHT { newDirection = LEFT }
		case Right: if oldDirection != LEFT { newDirection = RIGHT }
		}
	})

	dead = this.UpdateWithInput(newDirection)
	return
}

func (this *SnakeGame) Draw() {
	if this.Config.ClearScreen {
		this.screen.Clear()
	}

	for i := 0; i < len(this.snake); i++ {
		pos := this.snake[i].pos
		pos.X = pos.X * 2 + 1
		pos.Y = pos.Y + 1

		var prevDir,curDir Direction

		curDir = this.snake[i].direction
		if i != len(this.snake) - 1 {
			prevDir = this.snake[i + 1].direction
		} else {
			prevDir = this.tailDirection
		}

		this.screen.PutString(snakeCharset[curDir][prevDir], pos)
	}

	if this.fruit.X != -1 && this.fruit.Y != -1 {
		fruitPos := NewPos(this.fruit.X * 2 + 1, this.fruit.Y + 1)
		this.screen.PutChar('\u2299', fruitPos)
	}

	if this.Config.DrawBorder {
		this.screen.PutString(strings.Repeat("-", 2 * this.Config.Size.X + 2), NewPos(0, 0))
		this.screen.PutString(strings.Repeat("-", 2 * this.Config.Size.X + 2), NewPos(0, this.Config.Size.Y + 1))

		for y:= 1; y < this.Config.Size.Y + 1; y++ {
			this.screen.PutChar('|', NewPos(0, y))
			this.screen.PutChar('|', NewPos(2 * this.Config.Size.X + 1, y))
		}
	}
}
