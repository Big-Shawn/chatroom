package main

import (
	"fmt"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sys/unix"
)

var wscol = 30

func init() {
	ws, err := unix.IoctlGetWinsize(syscall.Stdout, unix.TIOCGWINSZ)
	if err != nil {
		panic(err)
	}
	wscol = int(ws.Col)
}

func renderbar(count, total int) {
	barwidth := wscol - len("Progress: 100% []")
	done := int(float64(barwidth) * float64(count) / float64(total))

	fmt.Printf("Progress: \x1b[33m%3d%%\x1b[0m ", count*100/total)
	fmt.Printf("[%s%s]",
		strings.Repeat("=", done),
		strings.Repeat("-", barwidth-done))
}

func main() {
	total := 50
	for i := 1; i <= total; i++ {
		//fmt.Print("\x1b7")   // save the cursor position
		fmt.Print("\033[s") // save the cursor position
		//fmt.Print("\x1b[2k") // erase the current line
		fmt.Print("\033[K") // erase the current line
		renderbar(i, total)
		time.Sleep(50 * time.Millisecond)
		//fmt.Print("\x1b8") // restore the cursor position
		fmt.Print("\033[u") // restore the cursor position

	}
	fmt.Println()
}
