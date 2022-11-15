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

// map[myDir][prevDir]
var snakeCharset = map[Direction]map[Direction]string {
	LEFT: {
		LEFT: "══",
		UP: "╗",
		DOWN: "╝",
	},
	RIGHT: {
		RIGHT: "══",
		UP: "╔═",
		DOWN: "╚═",
	},
	UP: {
		UP: "║",
		LEFT: "╚═",
		RIGHT: "╝",
	},
	DOWN: {
		DOWN: "║",
		LEFT: "╔═",
		RIGHT: "╗",
	},
}

type snakeSegment struct {
	pos Pos
	direction Direction
}

type SnakeGame struct {
	myRand *rand.Rand
	screen *Screen
	width, height uint
	delay time.Duration
	lastFrame time.Time
	snake []snakeSegment
	fruit Pos
}

func NewSnakeGame(screen *Screen, width, height uint, delay time.Duration) SnakeGame {
	source := rand.NewSource(time.Now().UnixNano())

	snake := make([]snakeSegment, 1)
	snake[0] = snakeSegment{NewPos(width / 2, height / 2), LEFT}

	fruit := NewPos(rand.Intn(int(width)), rand.Intn(int(height)))

	screen.StartInputReading()

	return SnakeGame{
		myRand: rand.New(source),
		screen: screen,
		width: width,
		height: height,
		delay: delay,
		lastFrame: time.Now(),
		snake: snake,
		fruit: fruit,
	}
}

func (this *SnakeGame) randomFruit() (fruit Pos) {
	outer:
	for {
		fruit = NewPos(rand.Intn(int(this.width)), rand.Intn(int(this.height)))
		for _, s := range this.snake {
			if fruit.Equal(s.pos) { continue outer }
		}
		return fruit
	}
}

func (this *SnakeGame) update() (dead bool) {
	head := this.snake[0]

	switch head.direction {
	case LEFT: head.pos.Addv(-1, 0)
	case RIGHT: head.pos.Addv(+1, 0)
	case UP: head.pos.Addv(0, -1)
	case DOWN: head.pos.Addv(0, +1)
	}

	if (head.pos.Equal(this.fruit)) {
		this.fruit = this.randomFruit()

		// If the snake has grown
		// we can simply insert the new head into the slice
		this.snake = append(this.snake, snakeSegment{NewPos(0, 0), head.direction})
		copy(this.snake[1:], this.snake[0:])
	} else {
		// Move the snake, by traslating all the positions
		// and discarding the last one
		if len(this.snake) > 1 {
			copy(this.snake[1:], this.snake[0:])
		}

		// Check if the snake collided with itself
		// This check is not done when the snake eats the fruit
		// because they cannot spawn into the snake
		for _, v := range this.snake {
			if v.pos.Equal(head.pos) {
				dead = true
				return
			}
		}
	}

	this.snake[0] = head

	dead = !head.pos.IsInside(NewPos(0, 0), NewPos(this.width - 1, this.height - 1))
	return
}

func (this *SnakeGame) draw() {
	this.screen.Clear()

	for i := 0; i < len(this.snake); i++ {
		pos := this.snake[i].pos
		pos.X = pos.X * 2 + 1
		pos.Y = pos.Y + 1

		var prevDir,curDir Direction

		curDir = this.snake[i].direction
		if i != len(this.snake) - 1 {
			prevDir = this.snake[i + 1].direction
		} else {
			prevDir = curDir
		}

		this.screen.PutString(snakeCharset[curDir][prevDir], pos)
	}

	fruitPos := NewPos(this.fruit.X * 2 + 1, this.fruit.Y + 1)

	this.screen.PutChar('\u2299', fruitPos)

	this.screen.PutString(strings.Repeat("-", int(2 * this.width + 2)), NewPos(0, 0))
	this.screen.PutString(strings.Repeat("-", int(2 * this.width + 2)), NewPos(0, this.height + 1))

	for y:= uint(1); y < this.width + 1; y++ {
		this.screen.PutChar('|', NewPos(0, y))
		this.screen.PutChar('|', NewPos(2 * this.width + 1, y))
	}

	this.lastFrame = this.screen.Draw(this.delay, this.lastFrame)
}

func (this *SnakeGame) RunGameLoop() {
	frameCount := 0
	for {
		this.draw()
		oldDirection := this.snake[0].direction
		direction := &this.snake[0].direction
		this.screen.ReadInput(func(k Keycode) {
			switch k {
			case Up: if oldDirection != DOWN { *direction = UP }
			case Down: if oldDirection != UP { *direction = DOWN }
			case Left: if oldDirection != RIGHT { *direction = LEFT }
			case Right: if oldDirection != LEFT { *direction = RIGHT }
			}
		})
		dead := this.update()
		if dead { break; }
		frameCount++
	}
}
