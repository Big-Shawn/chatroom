package main

import (
	"encoding/json"
	"fmt"
	"github.com/big-shawn/printer"
	"log"
	"net"
	"os"
	"os/signal"
	"reflect"
	"syscall"
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

/**
todo
1. 结束连接时需要上报服务器，清除连接信息
2. 定位屏幕输出，在指定位置输出内容 - done
3. 使用命令行工具进行工具初始化配置
4. 房间号，多个聊天室可同时开启
*/

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
	closeConn()

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
	//fmt.Println(string(msg))
	//fmt.Println(receive.Data)

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

func outputSendMsg(msg string) {
	timeString := time.Now().Format("2006-01-02 15:01:05")

	// 输出信息方消息
	sprintf := fmt.Sprintf("\r%s From Me : \n", timeString)
	printer.LastLine(sprintf)
	fmt.Printf("\r%s \n", string(msg))
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
		//fmt.Println("String Get", inputBuf)
		outputSendMsg(inputBuf)
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
	// 1. 用户主动退出  ctrl-c 退出
	osListener := make(chan os.Signal, 1)
	signal.Notify(osListener, syscall.SIGINT)
	// 2. 注册异常

}
