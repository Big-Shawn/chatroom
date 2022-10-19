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
var clientCon net.Conn

const MessageLimitSize = 1024

type ReceiveMsg struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}
type SendMsg struct {
	ClientId int    `json:"client_id"`
	Msg      []byte `json:"msg,omitempty"`
}

func main() {
	serverHost := "localhost"
	serverPort := "8088"
	closeSign = make(chan bool)

	dial, err := net.Dial("tcp", serverHost+":"+serverPort)
	if err != nil {
		panic(err)
	}
	clientCon = dial

	// receive Message
	go receiveMsg(dial)
	// send Message
	go getInput()

	// destruct Msg
	<-closeSign

}

func receiveMsg(con net.Conn) {
	for {
		buf := make([]byte, MessageLimitSize)
		n, err := con.Read(buf)
		if err != nil {
			panic(err)
		}
		outMsgResolve(buf[:n])
	}
}

// 信息输出
func outMsgResolve(msg []byte) {
	receive := &ReceiveMsg{}
	err := json.Unmarshal(msg, receive)
	fmt.Println(string(msg))
	fmt.Println(receive.Data)

	if err != nil {
		log.Println("Receive Error: " + err.Error())
		return
	}

	switch receive.Code {
	// clientId dispatch
	case 0:
		setClientId(receive)
	// client Msg
	case 1:
		outputClientMsg(receive)

	}

}

func outputClientMsg(msg *ReceiveMsg) {
	timeString := time.Now().Format("2006-01-02 15:01:05")
	reflectVal := reflect.ValueOf(msg.Data)
	from := int(reflectVal.MapIndex(reflect.ValueOf("from")).Interface().(float64))
	msgStr := reflectVal.MapIndex(reflect.ValueOf("msg")).Interface().(string)

	// 输出信息方消息
	fmt.Printf("\r%s From Client %d : \n", timeString, from)
	fmt.Printf("\r%s \n", msgStr)
}

// setClientId 设置由服务器分配下来的唯一性标识
func setClientId(msg *ReceiveMsg) {
	// 使用反射包进行接受数据的处理，从json到type
	data := msg.Data
	/**
	json decode 的时候会有默认的数据对应关系
	关系如下：
		bool, for JSON booleans
		float64, for JSON numbers
		string, for JSON strings
		[]interface{}, for JSON arrays
		map[string]interface{}, for JSON objects
		nil for JSON null
	*/
	clientId = int(
		reflect.ValueOf(data).MapIndex(
			reflect.ValueOf("clientId")).Interface().(float64))
}

func getInput() {
	for {
		//inputBuf := make([]byte, 1024)
		var inputBuf string
		_, err := fmt.Scanf("%s", &inputBuf)
		if err != nil {
			err := fmt.Errorf("input Error : %s", err.Error())
			fmt.Println(err)
			continue
		}
		fmt.Println("String Get", inputBuf)
		sendMsg([]byte(inputBuf))
	}
}

func sendMsg(input []byte) {
	//fmt.Println("String Get", string(input))
	//toSend := SendMsg{
	//	ClientId: clientId,
	//	Msg:      input,
	//}
	//marshal, err := json.Marshal(input)
	//if err != nil {
	//	err := fmt.Errorf("send msg Error: %s", err.Error())
	//	fmt.Println(err)
	//	return
	//}
	//
	//fmt.Println("Message Send Original: ", string(marshal))

	// 直接写原始数据 先不加多余的解释变量
	clientCon.Write(input)

	// 处理发送内容并携带自己的clientId
}

func closeConn() {

}
