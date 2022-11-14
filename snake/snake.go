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
	width, height uint
	delay time.Duration
	lastFrame time.Time
	direction Direction
	snake []Pos
	fruit Pos
}

func NewSnakeGame(screen *Screen, width, height uint, delay time.Duration) SnakeGame {
	source := rand.NewSource(time.Now().UnixNano())

	snake := make([]Pos, 1)
	snake[0] = Pos{X: int(width / 2), Y: int(height / 2)}

	fruit := NewPos(rand.Intn(int(width)), rand.Intn(int(height)))

	screen.StartInputReading()

	return SnakeGame{
		myRand: rand.New(source),
		screen: screen,
		width: width,
		height: height,
		delay: delay,
		lastFrame: time.Now(),
		direction: LEFT,
		snake: snake,
		fruit: fruit,
	}
}

func (this *SnakeGame) randomFruit() (fruit Pos) {
	outer:
	for {
		fruit = NewPos(rand.Intn(int(this.width)), rand.Intn(int(this.height)))
		for _, s := range this.snake {
			if fruit.Equal(s) { continue outer }
		}
		return fruit
	}
}

func (this *SnakeGame) update() (dead bool) {
	head := this.snake[0]

	switch this.direction {
	case LEFT: head.Addv(-1, 0)
	case RIGHT: head.Addv(+1, 0)
	case UP: head.Addv(0, -1)
	case DOWN: head.Addv(0, +1)
	}

	if (head.Equal(this.fruit)) {
		this.fruit = this.randomFruit()

		// If the snake has grown
		// we can simply insert the new head into the slice
		this.snake = append(this.snake, NewPos(0, 0))
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
			if v == head {
				dead = true
				return
			}
		}
	}

	this.snake[0] = head

	dead = !head.IsInside(NewPos(0, 0), NewPos(this.width - 1, this.height - 1))
	return
}

func (this *SnakeGame) draw() {
	this.screen.Clear()

	for _, s := range this.snake {
		// Shift for the borders
		s.X = s.X * 2 + 1
		s.Y = s.Y + 1

		this.screen.PutString("#", s)
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
		oldDirection := this.direction
		this.screen.ReadInput(func(k Keycode) {
			switch k {
			case Up: if oldDirection != DOWN { this.direction = UP }
			case Down: if oldDirection != UP { this.direction = DOWN }
			case Left: if oldDirection != RIGHT { this.direction = LEFT }
			case Right: if oldDirection != LEFT { this.direction = RIGHT }
			}
		})
		dead := this.update()
		if dead { break; }
		frameCount++
	}
}
