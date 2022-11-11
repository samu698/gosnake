package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
	"golang.org/x/sys/unix"
)

func clear() {
	fmt.Printf("\x1B[2J")
}

func setPos(x int, y int) {
	fmt.Printf("\x1B[%d;%dH", x + 1, y + 1)
}


type cellState int
const (
	EMPTY cellState = iota
	SNAKE
	FRUIT
)

type pos struct {
	x int
	y int
}

type direction int
const (
	LEFT direction = iota
	RIGHT
	UP
	DOWN
)


func randomFruit(width int, height int, snake []pos) pos {
	outer:
	for {
		p := pos{rand.Intn(width), rand.Intn(height)}
		for _, s := range snake {
			if p.x == s.x && p.y == s.y { continue outer }
		}
		return p
	}
}

func draw(width int, height int, snake []pos, fruit pos) {
	grid := make([]cellState, width * height)
	grid[fruit.y * width + fruit.x] = FRUIT
	for _, s := range snake {
		grid[s.y * width + s.x] = SNAKE
	}

	clear()	
	setPos(0, 0)
	fmt.Println(strings.Repeat("-", 2 * width + 1))
	for y := 0; y < height; y++ {
		fmt.Print("|")
		for x := 0; x < width - 1; x++ {
			switch grid[y * width + x] {
			case EMPTY: fmt.Print("  ")
			case SNAKE: fmt.Print("# ")
			case FRUIT: fmt.Print("\u2299 ")
			}
		}
		switch grid[y * width + width - 1] {
		case EMPTY: fmt.Print(" ")
		case SNAKE: fmt.Print("#")
		case FRUIT: fmt.Print("\u2299")
		}
		fmt.Println("|")
	}
	fmt.Println(strings.Repeat("-", 2 * width + 2))
}

func moveSnake(width int, height int, dir direction, snake *[]pos, fruit *pos) bool {
	var dx, dy int
	switch dir {
	case LEFT: dx = -1
	case RIGHT: dx = +1
	case UP: dy = -1
	case DOWN: dy = +1
	}

	hasEaten := fruit.x == (*snake)[0].x && fruit.y == (*snake)[0].y
	
	head := (*snake)[0]
	if len(*snake) > 1 {
		copy((*snake)[1:], (*snake)[0:])
	}

	tail := (*snake)[len(*snake) - 1]

	if (hasEaten) {
		*snake = append(*snake, tail)
		*fruit = randomFruit(width, height, *snake)
	}

	head.x += dx
	head.y += dy

	for _, v := range *snake {
		if v == head {
			return true
		}
	}

	(*snake)[0] = head

	return (*snake)[0].x < 0 || (*snake)[0].x >= width || (*snake)[0].y < 0 || (*snake)[0].y >= height
}

func openTTY() {
	info, err := unix.IoctlGetTermios(int(os.Stdin.Fd()), unix.TCGETS)
	if (err != nil) { fmt.Println(err); return }
	info.Lflag &^= (unix.ECHO | unix.ICANON)
	err = unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TCSETS, info)
	if (err != nil) { fmt.Println(err) }
}

func setBlocking() {
	opt, _ := unix.FcntlInt(os.Stdin.Fd(), unix.F_GETFL, 0)
	unix.FcntlInt(os.Stdin.Fd(), unix.F_SETFL, opt &^ unix.O_NONBLOCK)
}
func setNonBlocking() {
	opt, _ := unix.FcntlInt(os.Stdin.Fd(), unix.F_GETFL, 0)
	unix.FcntlInt(os.Stdin.Fd(), unix.F_SETFL, opt | unix.O_NONBLOCK)
}

func closeTTY() {
}


func menu(elements []string) int {
	var buf [16]byte
	for i, v := range elements {
		fmt.Printf("[%d] %s\n", i + 1, v)
	}
	os.Stdin.Read(buf[:])
	for i := range elements {
		if int(buf[0]) == int('1') + i {
			return i
		}
	}
	return 0
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var width, height int
	fmt.Print("insert width: ")
	fmt.Scan(&width)
	fmt.Print("insert height: ")
	fmt.Scan(&height)

	snake := make([]pos, 1)
	snake[0] = pos{width / 2, height / 2}

	dir := LEFT
	fruit := randomFruit(width, height, snake)

	openTTY()

	var buf [16]byte

	delay := time.Millisecond
	fmt.Println("Select speed: ")
	switch menu([]string{"SLOW", "MEDIUM", "FAST", "SONIC"}) {
	case 0: delay *= 120
	case 1: delay *= 95
	case 2: delay *= 70
	case 3: delay *= 50
	}

	setNonBlocking()

	for {
		os.Stdin.Read(buf[:])
		if buf[0] == 27 && buf[1] == 91 {
			switch buf[2] {
			case 65: if dir != DOWN { dir = UP }
			case 66: if dir != UP { dir = DOWN }
			case 68: if dir != RIGHT { dir = LEFT }
			case 67: if dir != LEFT { dir = RIGHT }
			}
		}

		if moveSnake(width, height, dir, &snake, &fruit) { return }
		draw(width, height, snake, fruit)
		time.Sleep(delay)
	}

}
