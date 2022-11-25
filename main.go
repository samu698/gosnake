package main

import (
	"fmt"
	"os"
	"time"
	. "github.com/samu698/gosnake/screen"
	. "github.com/samu698/gosnake/snake"
	. "github.com/samu698/gosnake/menu"
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

	/*
	puppeteer := newPuppeteer(&screen, paths)
	*/

	/*
	for i := 0; i < 20; i++ {
		puppeteer.Draw()
		puppeteer.Update()
		lastFrame = screen.Swap(delay, lastFrame)
	}
	*/

	playButton := NewButton("Play Snake!")
	frameDelay := NewIntEntry("Frame Delay: %dms", 100, 50, 500, 5)
	width := NewIntEntry("Width: %d cells", 40, 10, screen.GetSize().X / 2, 1)
	height := NewIntEntry("Height: %d cells", 40, 10, screen.GetSize().Y - 1, 1)

	menu := NewGameMenu(&screen, []MenuEntry{
		&playButton,
		&frameDelay,
		&width,
		&height,
	})
	
	lastFrame := time.Now()
	screen.StartInputReading()
	for {
		menu.Draw()
		menu.Update()
		if playButton.Pressed { break }
		lastFrame = screen.Swap(time.Millisecond * 33, lastFrame)
	}

	frameDelayMs := time.Duration(frameDelay.Value) * time.Millisecond

	config := SnakeConfig{
		Size: NewPos(width.Value, height.Value),
		Origin: NewPos(width.Value / 2, height.Value / 2),
		StartDirection: LEFT,
		CheckBounds: true,
		DrawBorder: true,
		SpawnFruit: true,
		GrowOnFrame: false,
		ClearScreen: true,
	}

	snakeGame := NewSnakeGame(&screen, config)
	snakeGame.Draw()
	for {
		if snakeGame.Update() { break }
		snakeGame.Draw()
		lastFrame = screen.Swap(frameDelayMs, lastFrame)
	}

	screen.Restore()
}
