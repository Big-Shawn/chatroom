package main

import (
	"fmt"
	"github.com/big-shawn/printer"
)

func main() {
	var line string
	_, err := fmt.Scan(&line)
	if err != nil {
		panic(err)
	}
	printer.LastLine(line + "========")
}
