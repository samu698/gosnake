package screen

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"golang.org/x/sys/unix"
)

type drawCmd struct {
	pos Pos
	char rune
}

type Screen struct {
	width, height uint
	termios unix.Termios
	fdFlags int

	drawQueue []drawCmd
	clearRequested bool

	inputBuffer []Keycode
	inputMutex sync.Mutex
	readerCancel context.CancelFunc
}

func NewScreen() Screen {
	winsize, _ := unix.IoctlGetWinsize(int(os.Stdin.Fd()), unix.TIOCGWINSZ)
	termios, _ := unix.IoctlGetTermios(int(os.Stdin.Fd()), unix.TCGETS)
	fdFlags, _ := unix.FcntlInt(os.Stdin.Fd(), unix.F_GETFL, 0)

	fmt.Print("\x1B[?1049h") // Enter alternate screen

	return Screen{
		width: uint(winsize.Col),
		height: uint(winsize.Row),
		termios: *termios,
		fdFlags: fdFlags,
		drawQueue: make([]drawCmd, 0),
		clearRequested: false,
		inputBuffer: make([]Keycode, 0),
		readerCancel: nil,
	}
}

func (this *Screen) GetSize() Pos {
	return NewPos(this.width, this.height)
}

func (this *Screen) Clear() {
	this.drawQueue = this.drawQueue[:0]
	this.clearRequested = true
}

func (this *Screen) PutChar(char rune, pos Pos) {
	if !pos.IsInside(NewPos(0, 0), NewPos(this.width, this.height)) {
		return
	}
	this.drawQueue = append(this.drawQueue, drawCmd{pos, char})
}

func (this *Screen) PutString(str string, pos Pos) {
	for _, char := range str {
		this.PutChar(char, pos)
		pos.X++
	}
}

func (this *Screen) Swap(frameTime time.Duration, lastFrame time.Time) time.Time {
	if this.clearRequested {
		fmt.Printf("\x1B[2J")
	}
	
	for _, cmd := range this.drawQueue {
		fmt.Printf("\x1B[%d;%dH%c", cmd.pos.Y + 1, cmd.pos.X + 1, cmd.char)
	}

	time.Sleep(time.Until(lastFrame.Add(frameTime)))
	return time.Now()
}

func (this *Screen) Restore() {
	fmt.Print("\x1B[?1049l") // Exit alternate screen
	this.StopInputReading()
}
