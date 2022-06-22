package selective

import (
	"errors"
	"os"
	"syscall"
	"time"
	"unsafe"

	"github.com/nsf/termbox-go"
)

const defaultRatio float64 = 7.0 / 3.0 // The terminal's default cursor width/height ratio

var (
	width int
	height int
	xwritten []int
	ywritten []int
	// whratio float64
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		if intInSlice(x, xwritten) && intInSlice(y, ywritten) {
			continue
		}
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func tbprintrev(x, y int, fg, bg termbox.Attribute, msg string) {
	x -= len(msg)
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func Select(options []string) (int, error) {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)

	selected := 0

	termbox.SetOutputMode(termbox.Output256)
	draw(options, selected)
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Key == termbox.KeyEsc || ev.Ch == 'q' || ev.Key == termbox.KeyCtrlC {
				return 0, errors.New("closed unexpectedly")
			}
			if ev.Key == termbox.KeyArrowRight && len(options)-1 > selected {
				selected++
				draw(options, selected)
			}
			if ev.Key == termbox.KeyArrowLeft && selected > 0 {
				selected--
				draw(options, selected)
			}
			if ev.Key == termbox.KeyEnter {
				return selected, nil
			}
		case termbox.EventResize:
			draw(options, selected)
		case termbox.EventMouse:

		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func draw(options []string, selected int) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	defer termbox.Flush()

	width, height, _ = canvasSize()

	for i, str := range options {
		if i == selected {
			tbprint(0, i, termbox.ColorWhite, termbox.ColorDefault, "> ")
			tbprint(2, i, termbox.ColorBlue, termbox.ColorDefault, str)
			continue
		}
		tbprint(2, i, termbox.ColorDefault, termbox.ColorDefault, str)
	}
	closeMsg := "To close, press the 'ESC' key"
	helpMsg := "Use the left and right arrow keys to select an option"
	tbprintrev(width-1, height-1, termbox.ColorWhite, termbox.ColorDefault, closeMsg)
	if width-2 > len(closeMsg) + len(helpMsg) { // '2' is subtracted for adequate spacing between sentences
		tbprint(0, height-1, termbox.ColorWhite, termbox.ColorDefault, helpMsg)
	} else {
		tbprint(0, height-2, termbox.ColorWhite, termbox.ColorDefault, helpMsg)
	}
}

// canvasSize returns the terminal columns, rows, and cursor aspect ratio
func canvasSize() (int, int, float64) {
	var size [4]uint16
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(os.Stdout.Fd()), uintptr(syscall.TIOCGWINSZ), uintptr(unsafe.Pointer(&size)), 0, 0, 0); err != 0 {
		panic(err)
	}
	rows, cols, width, height := size[0], size[1], size[2], size[3]

	var whratio = defaultRatio
	if width > 0 && height > 0 {
		whratio = float64(height/rows) / float64(width/cols)
	}

	return int(cols), int(rows), whratio
}

func intInSlice(a int, list []int) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
}