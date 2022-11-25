package menu

import (
	"fmt"

	. "github.com/samu698/gosnake/puppeteer"
	. "github.com/samu698/gosnake/screen"
	. "github.com/samu698/gosnake/snake"
)

var letterS = []Direction{RIGHT, RIGHT, RIGHT, RIGHT, UP, UP, LEFT, LEFT, LEFT, UP, UP, RIGHT, RIGHT}
var letterN = []Direction{UP, UP, UP, UP, UP, RIGHT, DOWN, RIGHT, DOWN, DOWN, RIGHT, DOWN, RIGHT, UP, UP, UP, UP}
var letterA1 = []Direction{UP, UP, UP, UP, UP, RIGHT, RIGHT, RIGHT, DOWN, DOWN, DOWN, DOWN}
var letterA2 = []Direction{LEFT, LEFT}
var letterK1 = []Direction{UP, UP, UP, UP, UP}
var letterK2 = []Direction{DOWN, LEFT, DOWN, LEFT, DOWN, RIGHT, DOWN, RIGHT, DOWN}
var letterE1 = []Direction{LEFT, LEFT, LEFT, LEFT, DOWN, DOWN, RIGHT, RIGHT, RIGHT}
var letterE2 = []Direction{LEFT, LEFT, LEFT, LEFT, UP}

var paths = []Path {
	{ NewPos(20, 20), letterS },
	{ NewPos(25, 20), letterN },
	{ NewPos(31, 20), letterA1 },
	{ NewPos(33, 18), letterA2 },
	{ NewPos(36, 20), letterK1 },
	{ NewPos(39, 16), letterK2 },
	{ NewPos(44, 16), letterE1 },
	{ NewPos(44, 20), letterE2 },
}

type MenuEntry interface {
	getText() string
	onKey(Keycode)
}

type ButtonMenuEntry struct {
	Text string
	Pressed bool
}

func NewButton(text string) ButtonMenuEntry {
	return ButtonMenuEntry{
		Text: text,
		Pressed: false,
	}
}

func (this *ButtonMenuEntry) getText() string {
	return this.Text
}

func (this *ButtonMenuEntry) onKey(key Keycode) {
	if (key == '\n') {
		this.Pressed = true
	}
}

type IntMenuEntry struct {
	// This must contain one "%d" format sequence
	FormatText string
	
	Value, Min, Max, Step int
}

func NewIntEntry(formatText string, value, min, max, step int) IntMenuEntry {
	return IntMenuEntry{
		FormatText: formatText,
		Value: value,
		Min: min,
		Max: max,
		Step: step,
	}
}

func (this *IntMenuEntry) getText() string {
	return fmt.Sprintf(this.FormatText, this.Value)
}

func (this *IntMenuEntry) onKey(key Keycode) {
	if (key == Right && this.Value + this.Step < this.Max) {
		this.Value += this.Step
	} else if (key == Left && this.Value - this.Step > this.Min) {
		this.Value -= this.Step
	}
}

type GameMenu struct {
	screen *Screen

	entries []MenuEntry
	selectedIndex int

	puppeteer Puppeteer
}

func NewGameMenu(screen *Screen, entries []MenuEntry) GameMenu {
	return GameMenu{
		screen: screen,
		entries: entries,
		selectedIndex: 0,
	}
}

func (this *GameMenu) Update() {
	this.screen.ReadInput(func (k Keycode) {
		if (k == Up && this.selectedIndex > 0) {
			this.selectedIndex--
		} else if (k == Down && this.selectedIndex < len(this.entries) - 1) {
			this.selectedIndex++
		} else {
			this.entries[this.selectedIndex].onKey(k)
		}
	})
}

func (this *GameMenu) Draw() {
	this.screen.Clear()
	for i, entry := range this.entries {
		entryStr := entry.getText()
		if (i == this.selectedIndex) {
			entryStr = "> " + entryStr + " <"
		}
		xPos := (this.screen.GetSize().X - len(entryStr)) / 2
		this.screen.PutString(entryStr, NewPos(xPos, i + 1))
	}
}
