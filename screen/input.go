package screen

import (
	"context"
	"os"
	"unicode/utf8"

	"golang.org/x/sys/unix"
)

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

func (this *Screen) inputReader(ctx context.Context) {
	var buf[32]byte
	for {
		select {
		case <-ctx.Done():
			return
		default:
			bytesRead, _ := os.Stdin.Read(buf[:])
			if bytesRead != 0 {
				keycode := keyCodeFromBytes(buf[:bytesRead])
				this.inputMutex.Lock()
				this.inputBuffer = append(this.inputBuffer, keycode)
				this.inputMutex.Unlock()
			}
		}
	}
}

func (this *Screen) StartInputReading() {
	if this.readerCancel != nil {
		return
	}

	// Disable echo and canonical mode
	// This way we can read the input as it comes, not line by line
	// Also the input won't be written on the terminal
	termios := this.termios
	termios.Lflag &^= unix.ECHO | unix.ICANON
	unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TCSETS, &termios)

	var ctx context.Context
	ctx, this.readerCancel = context.WithCancel(context.Background())

	go this.inputReader(ctx)
}

func (this *Screen) StopInputReading() {
	if this.readerCancel == nil {
		return
	}

	// Reset the flags to the original state
	unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TCSETS, &this.termios)

	this.readerCancel()
	this.readerCancel = nil
}


func (this *Screen) ReadInput(callback func(Keycode)) {
	this.inputMutex.Lock()
	for _, keycode := range this.inputBuffer {
		callback(keycode)
	}
	this.inputBuffer = this.inputBuffer[:0]
	this.inputMutex.Unlock()
}
