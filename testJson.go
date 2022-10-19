package main

import (
	"encoding/json"
	"fmt"
	"go/printer"
	"golang.org/x/sys/unix"
)

type Msg struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

func main() {
	printer.TerminalWindowSize()
	unix.Syscall()

	printer.Mode()
	var input string
	//input = make([]byte, 1024)
	_, err2 := fmt.Scan(&input)
	if err2 != nil {
		fmt.Println(err2)
	}

	m1 := Msg{
		Code: 0,
		Msg:  "Issue Client Id",
		Data: map[string]interface{}{
			"clientId": 1,
			"msg":      input,
		},
	}

	marshal, err := json.Marshal(m1)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(marshal))

	receive := &Msg{}
	err2 = json.Unmarshal(marshal, receive)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(receive.Code)

}
