package Printer

import (
	f "fmt"
	"golang.org/x/sys/unix"
	"syscall"
	"unsafe"
)

type Info struct {
	X      int // 当前光标
	Y      int
	Height int
	Weight int
}

func (i *Info) goesTo(x, y int) {
	f.Printf("\033[%d;%dH", y, x)
	i.X = x
	i.Y = y
}

// 清屏
func (i *Info) clear() {
	i.X = 0
	i.Y = 0
	f.Printf("\033[2J")
}

//保存光标位置
func savePosition() {
	f.Printf("\033[s")
}

//恢复光标位置
func restorePosition() {
	f.Printf("\033[u")
}

func moveUp(row int) {
	f.Printf("\033[%dA", row)
}

func moveDown(row int) {
	f.Printf("\033[%dB", row)
}

func moveLeft(row int) {
	f.Printf("\033[%dC", row)
}

func moveRight(row int) {
	f.Printf("\033[%dD", row)
}

func toHead() {
	f.Printf("\r")
}

func LastLine(line string) {
	moveUp(1)
	toHead()
	f.Printf(line)
}

// 获取窗口大小
type winsize struct {
	Row uint16
	Col uint16
	X   uint16
	Y   uint16
}
type window struct {
	Row uint16
	Col uint16
}

func TerminalWindowSize() (window, error) {
	unix.IoctlGetWinsize(syscall.Stdout, unix.TIOCGWINSZ)
	win := window{0, 0}
	res, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ), //此参数,不同的操作系统可能不一样,例如:TIOCGWINSZ_OSX
		uintptr(unsafe.Pointer(&win)),
	)
	if int(res) == -1 {
		return window{0, 0}, err
	}

	return win, nil
}

func getWinSize(fd int) (row, col uint16, err error) {
	var ws *winsize
	retCode, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL, uintptr(fd),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))
	if int(retCode) == -1 {
		panic(errno)
	}
	return ws.Row, ws.Col, nil
}
