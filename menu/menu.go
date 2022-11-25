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
	{ StartingPos: NewPos(20, 15), Path: letterS },
	{ StartingPos: NewPos(25, 15), Path: letterN },
	{ StartingPos: NewPos(31, 15), Path: letterA1 },
	{ StartingPos: NewPos(33, 13), Path: letterA2 },
	{ StartingPos: NewPos(36, 15), Path: letterK1 },
	{ StartingPos: NewPos(39, 11), Path: letterK2 },
	{ StartingPos: NewPos(44, 11), Path: letterE1 },
	{ StartingPos: NewPos(44, 15), Path: letterE2 },
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
		puppeteer: NewPuppeteer(screen, paths),
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
	this.puppeteer.Update()
}

func (this *GameMenu) Draw() {
	this.screen.Clear()
	for i, entry := range this.entries {
		entryStr := entry.getText()
		if (i == this.selectedIndex) {
			entryStr = "> " + entryStr + " <"
		}
		xPos := (this.screen.GetSize().X - len(entryStr)) / 2
		this.screen.PutString(entryStr, NewPos(xPos, i + 20))
	}
	this.puppeteer.Draw()
}
