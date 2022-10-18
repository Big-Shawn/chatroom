package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"time"
)

var closeSign chan bool
var clientId int

type ReceiveMsg returnMsg

func main() {
	serverHost := "localhost"
	serverPort := "8088"
	closeSign = make(chan bool)

	dial, err := net.Dial("tcp", serverHost+":"+serverPort)
	if err != nil {
		panic(err)
	}

	// receive Message
	go receiveMsg(dial)
	// send Message

	// destruct Msg
	<-closeSign

}

func receiveMsg(con net.Conn) {
	for {
		buf := make([]byte, MessageLimitSize)
		_, err := con.Read(buf)
		if err != nil {
			panic(err)
		}
		outMsgResolve(buf)
	}
}

// 信息输出
func outMsgResolve(msg []byte) {
	receive := &returnMsg{}
	err := json.Unmarshal(msg, receive)
	if err != nil {
		log.Println("Receive Error: " + err.Error())
	}
	switch receive.code {
	// clientId dispatch
	case 0:
		setClientId(receive)
	// client Msg
	case 1:
		outputClientMsg(receive)

	}

}

func outputClientMsg(msg *returnMsg) {
	timeString := time.Now().String()
	reflectVal := reflect.ValueOf(msg.data)
	from := int(reflectVal.MapIndex(reflect.ValueOf("from")).Int())
	msgStr := reflectVal.MapIndex(reflect.ValueOf("msg")).String()

	// 输出信息方消息
	fmt.Printf("\r%s From Client %d : \n", timeString, from)
	fmt.Printf("\r%s \n", msgStr)
}

// setClientId 设置由服务器分配下来的唯一性标识
func setClientId(msg *returnMsg) {
	// 使用反射包进行接受数据的处理，从json到type
	data := msg.data
	dataMapIter := reflect.ValueOf(data).MapRange()
	for dataMapIter.Next() {
		if key := dataMapIter.Key(); key.String() == "clientId" {
			clientId = int(dataMapIter.Value().Int())
			// 分配异常时关闭进程通道
			if clientId < 0 {
				close(closeSign)
			}
			break
		}
	}
}

func sendMsg() {
	// getInput
}

func closeConn() {

}
