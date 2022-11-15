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
	delay := time.Millisecond

	for {
		var tmp string
		fmt.Print("Do you want default settings [Y/N]: ")
		fmt.Scan(&tmp)
		if tmp == "Y" || tmp == "y" {
			fmt.Print("You have selected auto settings:\n")
			width = 40
			height = 40
			delay *= 120
			break
		} else if tmp == "N" || tmp == "n" {
			fmt.Print("You have selected manual settings:\n")
			fmt.Print("insert width: ")
			fmt.Scan(&width)
			fmt.Print("insert height: ")
			fmt.Scan(&height)

			delay = time.Millisecond
			fmt.Println("Select speed: ")
			switch menu([]string{"SLOW", "MEDIUM", "FAST", "SONIC"}) {
			case 0:
				delay *= 120
			case 1:
				delay *= 95
			case 2:
				delay *= 70
			case 3:
				delay *= 50
			}
			break
		} else {
			fmt.Print("You have insert a invalid input, please retry.\ns")
		}
	}

	snakeGame := NewSnakeGame(&screen, width, height, delay)
	snakeGame.RunGameLoop()
}
