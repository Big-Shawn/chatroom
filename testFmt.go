package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

func scanPromt() {
	fmt.Print("\f Input: ")
}

type window struct {
	Row uint16
	Col uint16
}

func terminalWindowSize() (window, error) {
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

func main() {

	res, _ := terminalWindowSize()
	fmt.Println(res)
	//
	//fmt.Printf("\033[%d;%dH", 100, 200)
	//fmt.Printf("hello")
	//for {
	//	//var input string
	//	//scanPromt()
	//	//_, err := fmt.Scan(&input)
	//	//if err != nil {
	//	//	panic(err)
	//	//}
	//	fmt.Printf("\033[%d;%dH", 100, 200)
	//	fmt.Printf("hello")
	//}

}
