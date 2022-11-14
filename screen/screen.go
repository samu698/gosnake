package screen

import (
	"fmt"
	"os"
	"time"
	"unicode/utf8"

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
}

type Keycode rune
const (
	Up Keycode = -(iota + 1)
	Down
	Right
	Left
	Unknown
)

func keyCodeFromBytes(bytes []byte) Keycode {
	sequences := map[string]Keycode {
		"\x1b[A": Up,
		"\x1b[B": Down,
		"\x1b[C": Right,
		"\x1b[D": Left,
	}

	if k, ok := sequences[string(bytes)]; ok {
		return k
	}

	r, _ := utf8.DecodeRune(bytes)
	if r != utf8.RuneError {
		return Keycode(r)
	}
	return Unknown
}

func NewScreen() Screen {
	winsize, _ := unix.IoctlGetWinsize(int(os.Stdin.Fd()), unix.TIOCGWINSZ)
	termios, _ := unix.IoctlGetTermios(int(os.Stdin.Fd()), unix.TCGETS)
	fdFlags, _ := unix.FcntlInt(os.Stdin.Fd(), unix.F_GETFL, 0)

	fmt.Fprintln(os.Stderr, winsize)

	return Screen{
		width: uint(winsize.Col),
		height: uint(winsize.Row),
		termios: *termios,
		fdFlags: fdFlags,
		drawQueue: make([]drawCmd, 0),
		clearRequested: false,
	}
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

// Disable echo and canonical mode
// This way we can read the input as it comes, not line by line
// Also the input won't be written on the terminal
func (this *Screen) SetDrawingFlags() {
	termios := this.termios
	termios.Lflag &^= unix.ECHO | unix.ICANON
	unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TCSETS, &termios)
}
func (this *Screen) ResetDrawingFlags() {
	unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TCSETS, &this.termios)
}

func (this *Screen) Draw(frameTime time.Duration, lastFrame time.Time) time.Time {
	if this.clearRequested {
		fmt.Printf("\x1B[2J")
	}
	
	for _, cmd := range this.drawQueue {
		fmt.Printf("\x1B[%d;%dH%c", cmd.pos.Y + 1, cmd.pos.X + 1, cmd.char)
	}

	time.Sleep(time.Until(lastFrame.Add(frameTime)))
	return time.Now()
}

func (this *Screen) ReadInput(callback func(Keycode)) {
	// Make stdin non blocking
	unix.FcntlInt(os.Stdin.Fd(), unix.F_SETFL, this.fdFlags | unix.O_NONBLOCK)

	var buf [32]byte
	bytesRead, _ := os.Stdin.Read(buf[:])
	if bytesRead != 0 {
		keycode := keyCodeFromBytes(buf[:bytesRead])
		callback(keycode)
	}

	// Restore stdin flags
	unix.FcntlInt(os.Stdin.Fd(), unix.F_SETFL, this.fdFlags)
}
