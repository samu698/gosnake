package main

import (
	"fmt"
	"os"
	"time"
	. "github.com/samu698/gosnake/screen"
	. "github.com/samu698/gosnake/snake"
)

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
	screen := NewScreen()

	var width, height uint
	fmt.Print("insert width: ")
	fmt.Scan(&width)
	fmt.Print("insert height: ")
	fmt.Scan(&height)

	delay := time.Millisecond
	fmt.Println("Select speed: ")
	switch menu([]string{"SLOW", "MEDIUM", "FAST", "SONIC"}) {
	case 0: delay *= 120
	case 1: delay *= 95
	case 2: delay *= 70
	case 3: delay *= 50
	}

	snakeGame := NewSnakeGame(&screen, width, height, delay)
	snakeGame.RunGameLoop()
}
